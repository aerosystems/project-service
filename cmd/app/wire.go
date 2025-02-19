//go:build wireinject
// +build wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/aerosystems/common-service/clients/gcpclient"
	"github.com/aerosystems/common-service/logger"
	"github.com/aerosystems/common-service/presenters/grpcserver"
	"github.com/aerosystems/common-service/presenters/httpserver"
	"github.com/aerosystems/project-service/internal/adapters"
	GRPCServer "github.com/aerosystems/project-service/internal/ports/grpc"
	HTTPServer "github.com/aerosystems/project-service/internal/ports/http"
	"github.com/aerosystems/project-service/internal/usecases"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(GRPCServer.ProjectUsecase), new(*usecases.ProjectUsecase)),
		wire.Bind(new(HTTPServer.ProjectUsecase), new(*usecases.ProjectUsecase)),
		wire.Bind(new(HTTPServer.TokenUsecase), new(*usecases.TokenUsecase)),
		wire.Bind(new(usecases.ProjectRepository), new(*adapters.ProjectRepo)),
		wire.Bind(new(usecases.SubscriptionAdapter), new(*adapters.SubscriptionAdapter)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideHTTPServer,
		ProvideLogrusLogger,
		ProvideHandler,
		ProvideProjectUsecase,
		ProvideTokenUsecase,
		ProvideProjectRepo,
		ProvideFirestoreClient,
		ProvideFirebaseAuthMiddleware,
		ProvideFirebaseAuthClient,
		ProvideGRPCServer,
		ProvideSubscriptionAdapter,
		ProvideProjectService,
	))
}

func ProvideApp(log *logrus.Logger, cfg *Config, httpServer *HTTPServer.Server, grpcServer *GRPCServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *Config {
	panic(wire.Build(NewConfig))
}

func ProvideHTTPServer(cfg *Config, log *logrus.Logger, firebaseAuth *HTTPServer.FirebaseAuth, handler *HTTPServer.Handler) *HTTPServer.Server {
	return HTTPServer.NewHTTPServer(&HTTPServer.Config{
		Config: httpserver.Config{
			Host: cfg.Host,
			Port: cfg.Port,
		},
		Mode: cfg.Mode,
	}, log, firebaseAuth, handler)
}

func ProvideGRPCServer(cfg *Config, log *logrus.Logger, projectService *GRPCServer.ProjectService) *GRPCServer.Server {
	return GRPCServer.NewGRPCServer(&grpcserver.Config{Host: cfg.Host, Port: cfg.Port}, log, projectService)
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideHandler(projectUsecase HTTPServer.ProjectUsecase, tokenUsecase HTTPServer.TokenUsecase) *HTTPServer.Handler {
	panic(wire.Build(HTTPServer.NewHandler))
}

func ProvideProjectUsecase(projectRepo usecases.ProjectRepository, subscriptionAdapter usecases.SubscriptionAdapter) *usecases.ProjectUsecase {
	panic(wire.Build(usecases.NewProjectUsecase))
}

func ProvideTokenUsecase(projectRepo usecases.ProjectRepository) *usecases.TokenUsecase {
	panic(wire.Build(usecases.NewTokenUsecase))
}

func ProvideFirestoreClient(cfg *Config) *firestore.Client {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GcpProjectId)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideProjectRepo(client *firestore.Client) *adapters.ProjectRepo {
	panic(wire.Build(adapters.NewProjectRepo))
}

func ProvideFirebaseAuthClient(cfg *Config) *auth.Client {
	client, err := gcpclient.NewFirebaseClient(cfg.GcpProjectId, cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return client
}

func ProvideFirebaseAuthMiddleware(client *auth.Client) *HTTPServer.FirebaseAuth {
	return HTTPServer.NewFirebaseAuth(client)
}

func ProvideSubscriptionAdapter(cfg *Config) *adapters.SubscriptionAdapter {
	subscriptionAdapter, err := adapters.NewSubscriptionAdapter(cfg.SubscriptionServiceGRPCAddr)
	if err != nil {
		panic(err)
	}
	return subscriptionAdapter
}

func ProvideProjectService(projectUsecase GRPCServer.ProjectUsecase) *GRPCServer.ProjectService {
	return GRPCServer.NewProjectService(projectUsecase)
}
