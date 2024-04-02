//go:build wireinject
// +build wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/aerosystems/project-service/internal/config"
	"github.com/aerosystems/project-service/internal/infrastructure/http"
	"github.com/aerosystems/project-service/internal/infrastructure/http/handlers"
	"github.com/aerosystems/project-service/internal/infrastructure/rpc"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/aerosystems/project-service/internal/repository/fire"
	"github.com/aerosystems/project-service/internal/repository/pg"
	"github.com/aerosystems/project-service/internal/repository/rpc"
	"github.com/aerosystems/project-service/internal/usecases"
	"github.com/aerosystems/project-service/pkg/gorm_postgres"
	"github.com/aerosystems/project-service/pkg/logger"
	"github.com/aerosystems/project-service/pkg/rpc_client"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(handlers.ProjectUsecase), new(*usecases.ProjectUsecase)),
		wire.Bind(new(handlers.TokenUsecase), new(*usecases.TokenUsecase)),
		wire.Bind(new(RpcServer.ProjectUsecase), new(*usecases.ProjectUsecase)),
		wire.Bind(new(usecases.SubsRepository), new(*RpcRepo.SubsRepo)),
		wire.Bind(new(usecases.ProjectRepository), new(*pg.ProjectRepo)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideHttpServer,
		ProvideRpcServer,
		ProvideLogrusLogger,
		ProvideLogrusEntry,
		ProvideGormPostgres,
		ProvideBaseHandler,
		ProvideProjectHandler,
		ProvideTokenHandler,
		ProvideProjectUsecase,
		ProvideTokenUsecase,
		ProvideSubsRepo,
		ProvideProjectRepo,
	))
}

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HttpServer.Server, rpcServer *RpcServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *config.Config {
	panic(wire.Build(config.NewConfig))
}

func ProvideHttpServer(log *logrus.Logger, cfg *config.Config, projectHandler *handlers.ProjectHandler, tokenHandler *handlers.TokenHandler) *HttpServer.Server {
	return HttpServer.NewServer(log, cfg.AccessSecret, projectHandler, tokenHandler)
}

func ProvideRpcServer(log *logrus.Logger, projectUsecase RpcServer.ProjectUsecase) *RpcServer.Server {
	panic(wire.Build(RpcServer.NewServer))
}

func ProvideLogrusEntry(log *logger.Logger) *logrus.Entry {
	return logrus.NewEntry(log.Logger)
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideGormPostgres(e *logrus.Entry, cfg *config.Config) *gorm.DB {
	db := GormPostgres.NewClient(e, cfg.PostgresDSN)
	if err := db.AutoMigrate(&models.Project{}); err != nil { // TODO: Move to migration
		panic(err)
	}
	return db
}

func ProvideBaseHandler(log *logrus.Logger, cfg *config.Config) *handlers.BaseHandler {
	return handlers.NewBaseHandler(log, cfg.Mode)
}

func ProvideProjectHandler(baseHandler *handlers.BaseHandler, projectUsecase handlers.ProjectUsecase) *handlers.ProjectHandler {
	panic(wire.Build(handlers.NewProjectHandler))
}

func ProvideTokenHandler(baseHandler *handlers.BaseHandler, tokenUsecase handlers.TokenUsecase) *handlers.TokenHandler {
	panic(wire.Build(handlers.NewTokenHandler))
}

func ProvideProjectUsecase(projectRepo usecases.ProjectRepository, subsRepo usecases.SubsRepository) *usecases.ProjectUsecase {
	panic(wire.Build(usecases.NewProjectUsecase))
}

func ProvideTokenUsecase(projectRepo usecases.ProjectRepository) *usecases.TokenUsecase {
	panic(wire.Build(usecases.NewTokenUsecase))
}

func ProvideSubsRepo(cfg *config.Config) *RpcRepo.SubsRepo {
	rpcClient := RpcClient.NewClient("tcp", cfg.SubsServiceRPCAddress)
	return RpcRepo.NewSubsRepo(rpcClient)
}

func ProvideProjectRepo(db *gorm.DB) *pg.ProjectRepo {
	panic(wire.Build(pg.NewProjectRepo))
}

func ProvideFirestoreClient(cfg *config.Config) *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GcpProjectId)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideFireProjectRepo(client *firestore.Client) *fire.ProjectRepo {
	panic(wire.Build(fire.NewProjectRepo))
}
