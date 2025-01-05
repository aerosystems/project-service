package HTTPServer

import (
	"context"
	"errors"
	"firebase.google.com/go/auth"
	"fmt"
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type ctxKey int

const (
	accessTokenClaimsContextKey ctxKey = iota
)

type accessTokenClaims struct {
	Uuid        string
	Email       string
	Role        string
	DisplayName string
}

type UserClaims struct {
	Uuid        uuid.UUID
	Email       string
	Role        models.Role
	DisplayName string
}

func (s *Server) setupMiddleware() {
	s.addLog(s.log)
	s.addCORS()
}

func (s *Server) addLog(log *logrus.Logger) {
	s.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, values middleware.RequestLoggerValues) error {
			log.WithFields(logrus.Fields{
				"URI":    values.URI,
				"status": values.Status,
			}).Info("request")

			return nil
		},
	}))
	s.echo.Use(middleware.Recover())
}

func (s *Server) addCORS() {
	DefaultCORSConfig := middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete, http.MethodOptions},
	}
	s.echo.Use(middleware.CORSWithConfig(DefaultCORSConfig))
}

type FirebaseAuth struct {
	client *auth.Client
}

func NewFirebaseAuth(client *auth.Client) *FirebaseAuth {
	return &FirebaseAuth{
		client: client,
	}
}

func (fa FirebaseAuth) RoleBased(roles ...models.Role) echo.MiddlewareFunc {
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

			//if !isAccess(roles, token.Claims["role"].(string)) {
			//	return echo.NewHTTPError(http.StatusForbidden, "access denied")
			//}

			ctx = context.WithValue(ctx, accessTokenClaimsContextKey, accessTokenClaims{
				Uuid:  token.UID,
				Email: token.Claims["email"].(string),
				Role:  "customer",
			})

			fmt.Printf("User %s with email %s\n", token.UID, token.Claims["email"].(string))

			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}

func NewUserClaims(claims accessTokenClaims) (UserClaims, error) {
	uuid, err := uuid.Parse(claims.Uuid)
	if err != nil {
		return UserClaims{}, err
	}
	role := models.RoleFromString(claims.Role)
	if role == models.UnknownRole {
		return UserClaims{}, CustomErrors.ErrUnknownUserRole
	}
	return UserClaims{
		Uuid:        uuid,
		Email:       claims.Email,
		Role:        role,
		DisplayName: claims.DisplayName,
	}, nil
}

func GetUserClaimsFromContext(ctx context.Context) (UserClaims, error) {
	claims, ok := ctx.Value(accessTokenClaimsContextKey).(accessTokenClaims)
	if !ok {
		return UserClaims{}, CustomErrors.ErrForbidden
	}
	user, err := NewUserClaims(claims)
	if err != nil {
		return UserClaims{}, CustomErrors.ErrForbidden
	}
	return user, nil
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

func getAuthHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) == 0 {
		return "", errors.New("missing Authorization header")
	}
	return authHeader, nil
}
