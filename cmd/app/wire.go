//go:build wireinject
// +build wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/auth"
	"github.com/aerosystems/project-service/internal/config"
	"github.com/aerosystems/project-service/internal/infrastructure/adapters/broker"
	"github.com/aerosystems/project-service/internal/infrastructure/adapters/rpc"
	"github.com/aerosystems/project-service/internal/infrastructure/repository/fire"
	"github.com/aerosystems/project-service/internal/presenters/http"
	"github.com/aerosystems/project-service/internal/presenters/http/handlers"
	"github.com/aerosystems/project-service/internal/presenters/http/handlers/project"
	"github.com/aerosystems/project-service/internal/presenters/http/handlers/token"
	"github.com/aerosystems/project-service/internal/presenters/http/middleware"
	"github.com/aerosystems/project-service/internal/presenters/rpc"
	"github.com/aerosystems/project-service/internal/usecases"
	"github.com/aerosystems/project-service/pkg/firebase"
	"github.com/aerosystems/project-service/pkg/logger"
	PubSub "github.com/aerosystems/project-service/pkg/pubsub"
	"github.com/aerosystems/project-service/pkg/rpc_client"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(usecases.CheckmailEventsAdapter), new(*broker.CheckmailEventsAdapter)),
		wire.Bind(new(handlers.ProjectUsecase), new(*usecases.ProjectUsecase)),
		wire.Bind(new(handlers.TokenUsecase), new(*usecases.TokenUsecase)),
		wire.Bind(new(RpcServer.ProjectUsecase), new(*usecases.ProjectUsecase)),
		wire.Bind(new(usecases.SubsRepository), new(*RpcRepo.SubsRepo)),
		wire.Bind(new(usecases.ProjectRepository), new(*fire.ProjectRepo)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideHttpServer,
		ProvideRpcServer,
		ProvideLogrusLogger,
		ProvideBaseHandler,
		ProvideProjectHandler,
		ProvideTokenHandler,
		ProvideProjectUsecase,
		ProvideTokenUsecase,
		ProvideSubsRepo,
		ProvideProjectRepo,
		ProvideFirestoreClient,
		ProvideFirebaseAuthMiddleware,
		ProvideFirebaseAuthClient,
		ProvidePubSubClient,
		ProvideCheckmailEventAdapter,
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

func ProvideHttpServer(log *logrus.Logger, cfg *config.Config, firebaseAuthMiddleware *middleware.FirebaseAuth, projectHandler *project.Handler, tokenHandler *token.Handler) *HttpServer.Server {
	return HttpServer.NewServer(log, firebaseAuthMiddleware, projectHandler, tokenHandler)
}

func ProvideRpcServer(log *logrus.Logger, projectUsecase RpcServer.ProjectUsecase) *RpcServer.Server {
	panic(wire.Build(RpcServer.NewServer))
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideBaseHandler(log *logrus.Logger, cfg *config.Config) *handlers.BaseHandler {
	return handlers.NewBaseHandler(log, cfg.Mode)
}

func ProvideProjectHandler(baseHandler *handlers.BaseHandler, projectUsecase handlers.ProjectUsecase) *project.Handler {
	panic(wire.Build(project.NewProjectHandler))
}

func ProvideTokenHandler(baseHandler *handlers.BaseHandler, tokenUsecase handlers.TokenUsecase) *token.Handler {
	panic(wire.Build(token.NewTokenHandler))
}

func ProvideProjectUsecase(projectRepo usecases.ProjectRepository, subsRepo usecases.SubsRepository, checkmailEventsAdapter usecases.CheckmailEventsAdapter) *usecases.ProjectUsecase {
	panic(wire.Build(usecases.NewProjectUsecase))
}

func ProvideTokenUsecase(projectRepo usecases.ProjectRepository) *usecases.TokenUsecase {
	panic(wire.Build(usecases.NewTokenUsecase))
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

func ProvideProjectRepo(client *firestore.Client) *fire.ProjectRepo {
	panic(wire.Build(fire.NewProjectRepo))
}

func ProvideFirebaseAuthMiddleware(client *auth.Client) *middleware.FirebaseAuth {
	return middleware.NewFirebaseAuth(client)
}

func ProvideFirebaseAuthClient(cfg *config.Config) *auth.Client {
	app, err := firebaseApp.NewApp(cfg.GcpProjectId, cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return app.Client
}

func ProvidePubSubClient(cfg *config.Config) *PubSub.Client {
	client, err := PubSub.NewClientWithAuth(cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideCheckmailEventAdapter(pubSubClient *PubSub.Client, cfg *config.Config) *broker.CheckmailEventsAdapter {
	return broker.NewCheckmailEventsAdapter(pubSubClient, cfg.CheckmailTopicId, cfg.CheckmailSubName, cfg.CheckmailCreateAccessEndpoint)
}
