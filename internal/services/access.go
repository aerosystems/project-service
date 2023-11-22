package services

import (
	"github.com/golang-jwt/jwt"
)

type AccessTokenClaims struct {
	AccessUuid string `json:"accessUuid"`
	UserUuid   string `json:"userUuid"`
	UserRole   string `json:"userRole"`
	Exp        int    `json:"exp"`
	jwt.StandardClaims
}

type TokenService interface {
	GetAccessSecret() string
	DecodeAccessToken(tokenString string) (*AccessTokenClaims, error)
}

type AccessTokenServiceImpl struct {
	accessSecret string
}

func NewAccessTokenServiceImpl(accessSecret string) *AccessTokenServiceImpl {
	return &AccessTokenServiceImpl{
		accessSecret: accessSecret,
	}
}

func (r *AccessTokenServiceImpl) GetAccessSecret() string {
	return r.accessSecret
}

func (r *AccessTokenServiceImpl) DecodeAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(r.accessSecret), nil
	})
	if claims, ok := token.Claims.(*AccessTokenClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}
