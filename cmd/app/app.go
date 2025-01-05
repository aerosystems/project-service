package main

import (
	"github.com/aerosystems/project-service/internal/common/config"
	GRPCServer "github.com/aerosystems/project-service/internal/presenters/grpc"
	"github.com/aerosystems/project-service/internal/presenters/http"
	"github.com/sirupsen/logrus"
)

type App struct {
	log        *logrus.Logger
	cfg        *config.Config
	httpServer *HTTPServer.Server
	grpcServer *GRPCServer.Server
}

func NewApp(
	log *logrus.Logger,
	cfg *config.Config,
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
