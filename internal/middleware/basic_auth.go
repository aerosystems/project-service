package middleware

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func BasicAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, password, ok := c.Request().BasicAuth()

		if !ok || !checkCredentials(username, password) {
			c.Response().Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
		}

		return next(c)
	}
}

func checkCredentials(username, password string) bool {
	validUsername := os.Getenv("BASIC_AUTH_DOCS_USERNAME")
	validPassword := os.Getenv("BASIC_AUTH_DOCS_PASSWORD")

	return username == validUsername && password == validPassword
}
