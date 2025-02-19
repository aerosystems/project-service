package main

import (
	GRPCServer "github.com/aerosystems/project-service/internal/ports/grpc"
	"github.com/aerosystems/project-service/internal/ports/http"
	"github.com/sirupsen/logrus"
)

type App struct {
	log        *logrus.Logger
	cfg        *Config
	httpServer *HTTPServer.Server
	grpcServer *GRPCServer.Server
}

func NewApp(
	log *logrus.Logger,
	cfg *Config,
	httpServer *HTTPServer.Server,
	grpcServer *GRPCServer.Server,
) *App {
	return &App{
		log:        log,
		cfg:        cfg,
		httpServer: httpServer,
		grpcServer: grpcServer,
	}
}
