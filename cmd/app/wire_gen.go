// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/aerosystems/project-service/internal/config"
	"github.com/aerosystems/project-service/internal/infrastructure/adapters/rpc"
	"github.com/aerosystems/project-service/internal/infrastructure/repository/fire"
	"github.com/aerosystems/project-service/internal/infrastructure/repository/pg"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/aerosystems/project-service/internal/presenters/http"
	"github.com/aerosystems/project-service/internal/presenters/http/handlers"
	"github.com/aerosystems/project-service/internal/presenters/rpc"
	"github.com/aerosystems/project-service/internal/usecases"
	"github.com/aerosystems/project-service/pkg/gorm_postgres"
	"github.com/aerosystems/project-service/pkg/logger"
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
	tokenUsecase := ProvideTokenUsecase(projectRepo)
	tokenHandler := ProvideTokenHandler(baseHandler, tokenUsecase)
	server := ProvideHttpServer(logrusLogger, config, projectHandler, tokenHandler)
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

func ProvideRpcServer(log *logrus.Logger, projectUsecase RpcServer.ProjectUsecase) *RpcServer.Server {
	server := RpcServer.NewServer(log, projectUsecase)
	return server
}

func ProvideProjectHandler(baseHandler *handlers.BaseHandler, projectUsecase handlers.ProjectUsecase) *handlers.ProjectHandler {
	projectHandler := handlers.NewProjectHandler(baseHandler, projectUsecase)
	return projectHandler
}

func ProvideTokenHandler(baseHandler *handlers.BaseHandler, tokenUsecase handlers.TokenUsecase) *handlers.TokenHandler {
	tokenHandler := handlers.NewTokenHandler(baseHandler, tokenUsecase)
	return tokenHandler
}

func ProvideProjectUsecase(projectRepo usecases.ProjectRepository, subsRepo usecases.SubsRepository) *usecases.ProjectUsecase {
	projectUsecase := usecases.NewProjectUsecase(projectRepo, subsRepo)
	return projectUsecase
}

func ProvideTokenUsecase(projectRepo usecases.ProjectRepository) *usecases.TokenUsecase {
	tokenUsecase := usecases.NewTokenUsecase(projectRepo)
	return tokenUsecase
}

func ProvideProjectRepo(db *gorm.DB) *pg.ProjectRepo {
	projectRepo := pg.NewProjectRepo(db)
	return projectRepo
}

func ProvideFireProjectRepo(client *firestore.Client) *fire.ProjectRepo {
	projectRepo := fire.NewProjectRepo(client)
	return projectRepo
}

// wire.go:

func ProvideHttpServer(log *logrus.Logger, cfg *config.Config, projectHandler *handlers.ProjectHandler, tokenHandler *handlers.TokenHandler) *HttpServer.Server {
	return HttpServer.NewServer(log, cfg.AccessSecret, projectHandler, tokenHandler)
}

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

func ProvideBaseHandler(log *logrus.Logger, cfg *config.Config) *handlers.BaseHandler {
	return handlers.NewBaseHandler(log, cfg.Mode)
}

func ProvideSubsRepo(cfg *config.Config) *RpcRepo.SubsRepo {
	rpcClient := RpcClient.NewClient("tcp", cfg.SubsServiceRPCAddress)
	return RpcRepo.NewSubsRepo(rpcClient)
}

func ProvideFirestoreClient(cfg *config.Config) *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GcpProjectId)
	if err != nil {
		panic(err)
	}
	return client
}
