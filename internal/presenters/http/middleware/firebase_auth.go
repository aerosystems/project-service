package middleware

import (
	"context"
	"firebase.google.com/go/auth"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
)

type FirebaseAuth struct {
	client *auth.Client
}

func NewFirebaseAuth(client *auth.Client) *FirebaseAuth {
	return &FirebaseAuth{
		client: client,
	}
}

func (fa FirebaseAuth) RoleBased(roles ...models.KindRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			jwt, err := getTokenFromHeader(c.Request())
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}

			token, err := fa.client.VerifyIDToken(ctx, jwt)
			if err != nil {
				log.Info(err)
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
			}

			//if !isAccess(roles, token.Claims["role"].(string)) {
			//	return echo.NewHTTPError(http.StatusForbidden, "access denied")
			//}

			ctx = context.WithValue(ctx, userContextKey, User{
				Uuid:  token.UID,
				Email: token.Claims["email"].(string),
			})

			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
