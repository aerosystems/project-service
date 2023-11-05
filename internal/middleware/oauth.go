package middleware

import (
	"errors"
	"github.com/aerosystems/customer-service/internal/services"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
	"time"
)

func AuthTokenMiddleware(roles []string) echo.MiddlewareFunc {
	AuthorizationConfig := echojwt.Config{
		SigningKey:     []byte(os.Getenv("ACCESS_SECRET")),
		ParseTokenFunc: parseToken,
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		},
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return AuthorizationConfig.ErrorHandler(c, errors.New("missing Authorization header"))
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
				return AuthorizationConfig.ErrorHandler(c, errors.New("invalid token format"))
			}

			token := tokenParts[1]

			accessTokenClaims, err := services.DecodeAccessToken(token)
			if err != nil {
				return AuthorizationConfig.ErrorHandler(c, err)
			}

			if int64(accessTokenClaims.Exp) < time.Now().Unix() {
				return AuthorizationConfig.ErrorHandler(c, errors.New("token expired"))
			}

			roleFound := false
			for _, role := range roles {
				if accessTokenClaims.UserRole == role {
					roleFound = true
					break
				}
			}

			if !roleFound {
				return echo.NewHTTPError(http.StatusForbidden, "access denied")
			}

			return next(c)
		}
	}
}

func parseToken(c echo.Context, auth string) (interface{}, error) {
	_ = c
	accessTokenClaims, err := services.DecodeAccessToken(auth)
	if err != nil {
		return nil, err
	}

	return accessTokenClaims, nil
}
