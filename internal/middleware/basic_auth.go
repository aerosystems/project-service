package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type BasicAuthMiddleware interface {
	BasicAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type BasicAuthMiddlewareImpl struct {
	username string
	password string
}

func NewBasicAuthMiddlewareImpl(username, password string) *BasicAuthMiddlewareImpl {
	return &BasicAuthMiddlewareImpl{
		username: username,
		password: password,
	}
}

func (b *BasicAuthMiddlewareImpl) BasicAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, password, ok := c.Request().BasicAuth()

		if !ok || !b.checkCredentials(username, password) {
			c.Response().Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}

		return next(c)
	}
}

func (b *BasicAuthMiddlewareImpl) checkCredentials(username, password string) bool {
	return username == b.username && password == b.password
}
