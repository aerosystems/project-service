package main

import (
	"github.com/aerosystems/project-service/internal/handlers"
)

type Config struct {
	baseHandler *handlers.BaseHandler
}

func NewConfig(baseHandler *handlers.BaseHandler) *Config {
	return &Config{
		baseHandler: baseHandler,
	}
}
