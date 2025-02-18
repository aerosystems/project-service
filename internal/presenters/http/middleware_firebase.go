package HTTPServer

import (
	"context"
	"errors"
	"firebase.google.com/go/v4/auth"
	"github.com/aerosystems/project-service/internal/entities"
	"github.com/go-logrusutil/logrusutil/logctx"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type ctxKey int

const (
	userContextKey ctxKey = iota

	errMessageForbidden    = "access denied"
	errMessageUnauthorized = "invalid token"
)

type User struct {
	UUID uuid.UUID
	Role entities.Role
}

type FirebaseAuth struct {
	client *auth.Client
}

func NewFirebaseAuth(client *auth.Client) *FirebaseAuth {
	return &FirebaseAuth{
		client: client,
	}
}

func (fa FirebaseAuth) RoleBasedAuth(roles ...entities.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			logger := logctx.From(ctx)
			jwt, err := getTokenFromHeader(c.Request())
			if err != nil {
				logger.Errorf("could not get access token from header: %v", err)
				return echo.NewHTTPError(http.StatusUnauthorized, errMessageUnauthorized)
			}

			token, err := fa.client.VerifyIDToken(ctx, jwt)
			if err != nil {
				logger.WithField("token", jwt).Errorf("could not verify access token: %v : %v. JWT: %s", token, err, jwt)
				return echo.NewHTTPError(http.StatusUnauthorized, errMessageUnauthorized)
			}

			userUUID, ok := token.Claims["user_uuid"].(string)
			if !ok {
				logger.WithField("token", jwt).Errorf("user_uuid claim not found in access token. JWT: %s", jwt)
				return echo.NewHTTPError(http.StatusUnauthorized, errMessageUnauthorized)
			}

			var user User
			user.UUID, err = uuid.Parse(userUUID)
			if err != nil {
				logger.WithField("token", jwt).Errorf("could not parse user_uuid claim as uuid: %v", err)
				return echo.NewHTTPError(http.StatusUnauthorized, errMessageUnauthorized)
			}

			userRole, ok := token.Claims["role"].(string)
			if !ok {
				logger.WithField("token", jwt).Errorf("role claim not found in access token. JWT: %s", jwt)
				return echo.NewHTTPError(http.StatusForbidden, errMessageForbidden)
			}

			user.Role = entities.RoleFromString(userRole)
			if !isAccess(roles, user.Role) {
				logger.WithField("token", jwt).Errorf("user role %s is not allowed to access. JWT: %s", user.Role, jwt)
				return echo.NewHTTPError(http.StatusForbidden, errMessageForbidden)
			}

			ctx = context.WithValue(ctx, userContextKey, user)

			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func GetUserFromContext(ctx context.Context) (User, error) {
	user, ok := ctx.Value(userContextKey).(User)
	if !ok {
		return User{}, errors.New("user not found in context")
	}
	return user, nil
}

func isAccess(roles []entities.Role, role entities.Role) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func getAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return "", errors.New("missing Authorization header")
	}
	return authHeader, nil
}

func getTokenFromHeader(r *http.Request) (string, error) {
	authHeader, err := getAuthHeader(r)
	if err != nil {
		return "", err
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("invalid token format")
	}
	return tokenParts[1], nil
}
