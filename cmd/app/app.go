package main

import (
	"github.com/aerosystems/project-service/internal/config"
	"github.com/aerosystems/project-service/internal/http"
	"github.com/aerosystems/project-service/internal/infrastructure/rpc"
	"github.com/sirupsen/logrus"
)

type App struct {
	log        *logrus.Logger
	cfg        *config.Config
	httpServer *HTTPServer.Server
	rpcServer  *RPCServer.Server
}

func NewApp(
	log *logrus.Logger,
	cfg *config.Config,
	httpServer *HTTPServer.Server,
	rpcServer *RPCServer.Server,
) *App {
	return &App{
		log:        log,
		cfg:        cfg,
		httpServer: httpServer,
		rpcServer:  rpcServer,
	}
}
