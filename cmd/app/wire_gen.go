// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/aerosystems/project-service/internal/config"
	"github.com/aerosystems/project-service/internal/http"
	"github.com/aerosystems/project-service/internal/infrastructure/rest"
	"github.com/aerosystems/project-service/internal/infrastructure/rpc"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/aerosystems/project-service/internal/repository/pg"
	"github.com/aerosystems/project-service/internal/repository/rpc"
	"github.com/aerosystems/project-service/internal/usecases"
	"github.com/aerosystems/project-service/pkg/gorm_postgres"
	"github.com/aerosystems/project-service/pkg/logger"
	"github.com/aerosystems/project-service/pkg/oauth"
	"github.com/aerosystems/project-service/pkg/rpc_client"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Injectors from wire.go:

//go:generate wire
func InitApp() *App {
	logger := ProvideLogger()
	logrusLogger := ProvideLogrusLogger(logger)
	config := ProvideConfig()
	baseHandler := ProvideBaseHandler(logrusLogger, config)
	entry := ProvideLogrusEntry(logger)
	db := ProvideGormPostgres(entry, config)
	projectRepo := ProvideProjectRepo(db)
	subsRepo := ProvideSubsRepo(config)
	projectUsecase := ProvideProjectUsecase(projectRepo, subsRepo)
	projectHandler := ProvideProjectHandler(baseHandler, projectUsecase)
	tokenHandler := ProvideTokenHandler(baseHandler, projectUsecase)
	accessTokenService := ProvideAccessTokenService(config)
	server := ProvideHttpServer(logrusLogger, config, projectHandler, tokenHandler, accessTokenService)
	rpcServerServer := ProvideRpcServer(logrusLogger, projectUsecase)
	app := ProvideApp(logrusLogger, config, server, rpcServerServer)
	return app
}

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HttpServer.Server, rpcServer *RpcServer.Server) *App {
	app := NewApp(log, cfg, httpServer, rpcServer)
	return app
}

func ProvideLogger() *logger.Logger {
	loggerLogger := logger.NewLogger()
	return loggerLogger
}

func ProvideConfig() *config.Config {
	configConfig := config.NewConfig()
	return configConfig
}

func ProvideHttpServer(log *logrus.Logger, cfg *config.Config, projectHandler *rest.ProjectHandler, tokenHandler *rest.TokenHandler, tokenService HttpServer.TokenService) *HttpServer.Server {
	server := HttpServer.NewServer(log, projectHandler, tokenHandler, tokenService)
	return server
}

func ProvideRpcServer(log *logrus.Logger, projectUsecase RpcServer.ProjectUsecase) *RpcServer.Server {
	server := RpcServer.NewServer(log, projectUsecase)
	return server
}

func ProvideProjectHandler(baseHandler *rest.BaseHandler, projectUsecase rest.ProjectUsecase) *rest.ProjectHandler {
	projectHandler := rest.NewProjectHandler(baseHandler, projectUsecase)
	return projectHandler
}

func ProvideTokenHandler(baseHandler *rest.BaseHandler, projectUsecase rest.ProjectUsecase) *rest.TokenHandler {
	tokenHandler := rest.NewTokenHandler(baseHandler, projectUsecase)
	return tokenHandler
}

func ProvideProjectUsecase(projectRepo usecases.ProjectRepository, subsRepo usecases.SubsRepository) *usecases.ProjectUsecase {
	projectUsecase := usecases.NewProjectUsecase(projectRepo, subsRepo)
	return projectUsecase
}

func ProvideProjectRepo(db *gorm.DB) *pg.ProjectRepo {
	projectRepo := pg.NewProjectRepo(db)
	return projectRepo
}

// wire.go:

func ProvideLogrusEntry(log *logger.Logger) *logrus.Entry {
	return logrus.NewEntry(log.Logger)
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideGormPostgres(e *logrus.Entry, cfg *config.Config) *gorm.DB {
	db := GormPostgres.NewClient(e, cfg.PostgresDSN)
	if err := db.AutoMigrate(&models.Project{}); err != nil {
		panic(err)
	}
	return db
}

func ProvideBaseHandler(log *logrus.Logger, cfg *config.Config) *rest.BaseHandler {
	return rest.NewBaseHandler(log, cfg.Mode)
}

func ProvideSubsRepo(cfg *config.Config) *RpcRepo.SubsRepo {
	rpcClient := RPCClient.NewClient("tcp", cfg.SubsServiceRPCAddress)
	return RpcRepo.NewSubsRepo(rpcClient)
}

func ProvideAccessTokenService(cfg *config.Config) *OAuthService.AccessTokenService {
	return OAuthService.NewAccessTokenService(cfg.AccessSecret)
}
