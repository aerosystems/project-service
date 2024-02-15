package main

import (
	"github.com/aerosystems/project-service/internal/middleware"
)

type Config struct {
	baseHandler         *rest.BaseHandler
	oauthMiddleware     middleware.OAuthMiddleware
	basicAuthMiddleware middleware.BasicAuthMiddleware
}

func NewConfig(baseHandler *rest.BaseHandler, oauthMiddleware middleware.OAuthMiddleware, basicAuthMiddleware middleware.BasicAuthMiddleware) *Config {
	return &Config{
		baseHandler:         baseHandler,
		oauthMiddleware:     oauthMiddleware,
		basicAuthMiddleware: basicAuthMiddleware,
	}
}
