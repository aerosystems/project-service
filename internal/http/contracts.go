package HttpServer

import OAuthService "github.com/aerosystems/project-service/pkg/oauth"

type TokenService interface {
	GetAccessSecret() string
	DecodeAccessToken(tokenString string) (*OAuthService.AccessTokenClaims, error)
}
