package HTTPServer

import (
	"context"
	"errors"
	"firebase.google.com/go/auth"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type ctxKey int

const (
	userContextKey ctxKey = iota
)

type User struct {
	UUID  uuid.UUID
	Role  models.Role
	Email string
}

type FirebaseAuth struct {
	client *auth.Client
}

func NewFirebaseAuth(client *auth.Client) *FirebaseAuth {
	return &FirebaseAuth{
		client: client,
	}
}

func (fa FirebaseAuth) RoleBasedAuth(roles ...models.Role) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			jwt, err := getTokenFromHeader(c.Request())
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			token, err := fa.client.VerifyIDToken(ctx, jwt)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			userUUID, ok := token.Claims["user_uuid"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			var user User
			user.UUID, err = uuid.Parse(userUUID)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			userRole, ok := token.Claims["role"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			user.Role = models.RoleFromString(userRole)
			if !isAccess(roles, user.Role) {
				return echo.NewHTTPError(http.StatusForbidden, "access denied")
			}

			userEmail, ok := token.Claims["email"].(string)
			if !ok {
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			user.Email = userEmail

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

func isAccess(roles []models.Role, role models.Role) bool {
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
