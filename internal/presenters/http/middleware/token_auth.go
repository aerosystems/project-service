package middleware

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/golang-jwt/jwt"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"net/http"
)

type TokenBasedAuth struct {
	accessSecret string
}

func NewTokenBasedAuth(accessSecret string) *TokenBasedAuth {
	return &TokenBasedAuth{
		accessSecret: accessSecret,
	}
}

type AccessTokenClaims struct {
	AccessUuid string `json:"accessUuid"`
	UserUuid   string `json:"userUuid"`
	UserRole   string `json:"userRole"`
	Exp        int    `json:"exp"`
	jwt.StandardClaims
}

func (ta TokenBasedAuth) AuthTokenMiddleware(roles ...models.KindRole) echo.MiddlewareFunc {
	AuthorizationConfig := echojwt.Config{
		SigningKey:     []byte(ta.accessSecret),
		ParseTokenFunc: ta.parseToken,
		ErrorHandler: func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		},
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token, err := getTokenFromHeader(c.Request())
			if err != nil {
				return AuthorizationConfig.ErrorHandler(c, err)
			}
			accessTokenClaims, err := ta.DecodeAccessToken(token)
			if err != nil {
				return AuthorizationConfig.ErrorHandler(c, err)
			}
			if !isAccess(roles, accessTokenClaims.UserRole) {
				return echo.NewHTTPError(http.StatusForbidden, "access denied")
			}
			echo.Context(c).Set("accessTokenClaims", *accessTokenClaims)
			return next(c)
		}
	}
}

func (ta TokenBasedAuth) parseToken(c echo.Context, auth string) (interface{}, error) {
	_ = c
	accessTokenClaims, err := ta.DecodeAccessToken(auth)
	if err != nil {
		return nil, err
	}
	return accessTokenClaims, nil
}

func (ta TokenBasedAuth) DecodeAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(ta.accessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
