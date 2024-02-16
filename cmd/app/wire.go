//go:build wireinject
// +build wireinject

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
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(rest.ProjectUsecase), new(*usecases.ProjectUsecase)),
		wire.Bind(new(RPCServer.ProjectUsecase), new(*usecases.ProjectUsecase)),
		wire.Bind(new(usecases.SubsRepository), new(*rpcRepo.SubsRepo)),
		wire.Bind(new(usecases.ProjectRepository), new(*pg.ProjectRepo)),
		wire.Bind(new(HTTPServer.TokenService), new(*OAuthService.AccessTokenService)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideHTTPServer,
		ProvideRPCServer,
		ProvideLogrusLogger,
		ProvideLogrusEntry,
		ProvideGormPostgres,
		ProvideBaseHandler,
		ProvideProjectHandler,
		ProvideTokenHandler,
		ProvideProjectUsecase,
		ProvideSubsRepo,
		ProvideProjectRepo,
		ProvideAccessTokenService,
	))
}

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HTTPServer.Server, rpcServer *RPCServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *config.Config {
	panic(wire.Build(config.NewConfig))
}

func ProvideHTTPServer(log *logrus.Logger, cfg *config.Config, projectHandler *rest.ProjectHandler, tokenHandler *rest.TokenHandler, tokenService HTTPServer.TokenService) *HTTPServer.Server {
	panic(wire.Build(HTTPServer.NewServer))
}

func ProvideRPCServer(log *logrus.Logger, projectUsecase RPCServer.ProjectUsecase) *RPCServer.Server {
	panic(wire.Build(RPCServer.NewServer))
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

func ProvideBaseHandler(log *logrus.Logger, cfg *config.Config) *rest.BaseHandler {
	return rest.NewBaseHandler(log, cfg.Mode)
}

func ProvideProjectHandler(baseHandler *rest.BaseHandler, projectUsecase rest.ProjectUsecase) *rest.ProjectHandler {
	panic(wire.Build(rest.NewProjectHandler))
}

func ProvideTokenHandler(baseHandler *rest.BaseHandler, projectUsecase rest.ProjectUsecase) *rest.TokenHandler {
	panic(wire.Build(rest.NewTokenHandler))
}

func ProvideProjectUsecase(projectRepo usecases.ProjectRepository, subsRepo usecases.SubsRepository) *usecases.ProjectUsecase {
	panic(wire.Build(usecases.NewProjectUsecase))
}

func ProvideSubsRepo(cfg *config.Config) *rpcRepo.SubsRepo {
	rpcClient := RPCClient.NewClient("tcp", cfg.SubsServiceRPCAddress)
	return rpcRepo.NewSubsRepo(rpcClient)
}

func ProvideProjectRepo(db *gorm.DB) *pg.ProjectRepo {
	panic(wire.Build(pg.NewProjectRepo))
}

func ProvideAccessTokenService(cfg *config.Config) *OAuthService.AccessTokenService {
	return OAuthService.NewAccessTokenService(cfg.AccessSecret)
}
