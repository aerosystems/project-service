package main

import (
	"github.com/aerosystems/project-service/internal/handlers"
	"github.com/aerosystems/project-service/internal/middleware"
)

type Config struct {
	baseHandler         *handlers.BaseHandler
	oauthMiddleware     middleware.OAuthMiddleware
	basicAuthMiddleware middleware.BasicAuthMiddleware
}

func NewConfig(baseHandler *handlers.BaseHandler, oauthMiddleware middleware.OAuthMiddleware, basicAuthMiddleware middleware.BasicAuthMiddleware) *Config {
	return &Config{
		baseHandler:         baseHandler,
		oauthMiddleware:     oauthMiddleware,
		basicAuthMiddleware: basicAuthMiddleware,
	}
}
