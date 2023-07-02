package TokenService

import (
	"github.com/golang-jwt/jwt"
	"os"
)

type AccessTokenClaims struct {
	AccessUUID string `json:"accessUuid"`
	UserID     int    `json:"userId"`
	UserRole   string `json:"userRole"`
	Exp        int    `json:"exp"`
	jwt.StandardClaims
}

func DecodeAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})

	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
