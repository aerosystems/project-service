package models

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type AccessTokenClaims struct {
	AccessUuid string `json:"accessUuid"`
	UserUuid   string `json:"userUuid"`
	UserRole   string `json:"userRole"`
	Exp        int    `json:"exp"`
	jwt.StandardClaims
}

type RefreshTokenClaims struct {
	RefreshUuid string `json:"refreshUuid"`
	UserUuid    string `json:"userUuid"`
	UserRole    string `json:"userRole"`
	Exp         int    `json:"exp"`
	jwt.StandardClaims
}

// TokenDetails is the structure which holds data with JWT tokens
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUuid   uuid.UUID
	RefreshUuid  uuid.UUID
	AtExpires    int64
	RtExpires    int64
}
