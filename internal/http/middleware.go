package HTTPServer

import (
	"errors"
	"github.com/aerosystems/project-service/internal/models"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func (s *Server) setupMiddleware() {
	s.addLog(s.log)
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

func (s *Server) AuthTokenMiddleware(roles ...models.KindRole) echo.MiddlewareFunc {
	AuthorizationConfig := echojwt.Config{
		SigningKey:     []byte(s.tokenService.GetAccessSecret()),
		ParseTokenFunc: s.parseToken,
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		},
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := getTokenFromHeader(c)
			if err != nil {
				return AuthorizationConfig.ErrorHandler(c, err)
			}
			accessTokenClaims, err := s.tokenService.DecodeAccessToken(token)
			if err != nil {
				return AuthorizationConfig.ErrorHandler(c, err)
			}
			if !isAccess(roles, accessTokenClaims.UserRole) {
				return echo.NewHTTPError(http.StatusForbidden, "access denied")
			}
			echo.Context(c).Set("accessTokenClaims", accessTokenClaims)
			return next(c)
		}
	}
}

func (s *Server) parseToken(c echo.Context, auth string) (interface{}, error) {
	_ = c
	accessTokenClaims, err := s.tokenService.DecodeAccessToken(auth)
	if err != nil {
		return nil, err
	}
	return accessTokenClaims, nil
}

func isAccess(roles []models.KindRole, role string) bool {
	for _, r := range roles {
		if r.String() == role {
			return true
		}
	}
	return false
}

func getTokenFromHeader(c echo.Context) (string, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if len(authHeader) == 0 {
		return "", errors.New("missing Authorization header")
	}
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return "", errors.New("invalid token format")
	}
	return tokenParts[1], nil
}
