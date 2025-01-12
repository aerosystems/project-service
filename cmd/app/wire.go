//go:build wireinject
// +build wireinject

package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"firebase.google.com/go/v4/auth"
	"github.com/aerosystems/project-service/internal/adapters"
	"github.com/aerosystems/project-service/internal/common/config"
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/aerosystems/project-service/internal/common/protobuf/project"
	GRPCServer "github.com/aerosystems/project-service/internal/presenters/grpc"
	HTTPServer "github.com/aerosystems/project-service/internal/presenters/http"
	"github.com/aerosystems/project-service/internal/usecases"
	"github.com/aerosystems/project-service/pkg/gcp"
	"github.com/aerosystems/project-service/pkg/logger"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

//go:generate wire
func InitApp() *App {
	panic(wire.Build(
		wire.Bind(new(GRPCServer.ProjectUsecase), new(*usecases.ProjectUsecase)),
		wire.Bind(new(HTTPServer.ProjectUsecase), new(*usecases.ProjectUsecase)),
		wire.Bind(new(HTTPServer.TokenUsecase), new(*usecases.TokenUsecase)),
		wire.Bind(new(project.ProjectServiceServer), new(*GRPCServer.ProjectHandler)),
		wire.Bind(new(usecases.ProjectRepository), new(*adapters.ProjectRepo)),
		wire.Bind(new(usecases.SubscriptionAdapter), new(*adapters.SubscriptionAdapter)),
		ProvideApp,
		ProvideLogger,
		ProvideConfig,
		ProvideHTTPServer,
		ProvideLogrusLogger,
		ProvideBaseHandler,
		ProvideProjectHandler,
		ProvideTokenHandler,
		ProvideProjectUsecase,
		ProvideTokenUsecase,
		ProvideProjectRepo,
		ProvideFirestoreClient,
		ProvideFirebaseAuthMiddleware,
		ProvideFirebaseAuthClient,
		ProvideErrorHandler,
		ProvideGRPCHandlers,
		ProvideGRPCServer,
		ProvideSubscriptionAdapter,
	))
}

func ProvideApp(log *logrus.Logger, cfg *config.Config, httpServer *HTTPServer.Server, grpcServer *GRPCServer.Server) *App {
	panic(wire.Build(NewApp))
}

func ProvideLogger() *logger.Logger {
	panic(wire.Build(logger.NewLogger))
}

func ProvideConfig() *config.Config {
	panic(wire.Build(config.NewConfig))
}

func ProvideHTTPServer(cfg *config.Config, log *logrus.Logger, errorHandler *echo.HTTPErrorHandler, firebaseAuthMiddleware *HTTPServer.FirebaseAuth, projectHandler *HTTPServer.ProjectHandler, tokenHandler *HTTPServer.TokenHandler) *HTTPServer.Server {
	return HTTPServer.NewServer(cfg.Port, log, errorHandler, firebaseAuthMiddleware, projectHandler, tokenHandler)
}

func ProvideGRPCServer(cfg *config.Config, log *logrus.Logger, grpcHandler project.ProjectServiceServer) *GRPCServer.Server {
	return GRPCServer.NewGRPCServer(cfg.Port, log, grpcHandler)
}

func ProvideGRPCHandlers(projectUsecase GRPCServer.ProjectUsecase) *GRPCServer.ProjectHandler {
	panic(wire.Build(GRPCServer.NewProjectHandler))
}

func ProvideLogrusLogger(log *logger.Logger) *logrus.Logger {
	return log.Logger
}

func ProvideBaseHandler(log *logrus.Logger, cfg *config.Config) *HTTPServer.BaseHandler {
	return HTTPServer.NewBaseHandler(log, cfg.Mode)
}

func ProvideProjectHandler(baseHandler *HTTPServer.BaseHandler, projectUsecase HTTPServer.ProjectUsecase) *HTTPServer.ProjectHandler {
	panic(wire.Build(HTTPServer.NewProjectHandler))
}

func ProvideTokenHandler(baseHandler *HTTPServer.BaseHandler, tokenUsecase HTTPServer.TokenUsecase) *HTTPServer.TokenHandler {
	panic(wire.Build(HTTPServer.NewTokenHandler))
}

func ProvideProjectUsecase(projectRepo usecases.ProjectRepository, subscriptionAdapter usecases.SubscriptionAdapter) *usecases.ProjectUsecase {
	panic(wire.Build(usecases.NewProjectUsecase))
}

func ProvideTokenUsecase(projectRepo usecases.ProjectRepository) *usecases.TokenUsecase {
	panic(wire.Build(usecases.NewTokenUsecase))
}

func ProvideFirestoreClient(cfg *config.Config) *firestore.Client {
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

func ProvideFirebaseAuthMiddleware(client *auth.Client) *HTTPServer.FirebaseAuth {
	return HTTPServer.NewFirebaseAuth(client)
}

func ProvideFirebaseAuthClient(cfg *config.Config) *auth.Client {
	app, err := gcp.NewFirebaseApp(cfg.GcpProjectId, cfg.GoogleApplicationCredentials)
	if err != nil {
		panic(err)
	}
	return app.Client
}

func ProvideErrorHandler(cfg *config.Config) *echo.HTTPErrorHandler {
	errorHandler := CustomErrors.NewEchoErrorHandler(cfg.Mode)
	return &errorHandler
}

func ProvideSubscriptionAdapter(cfg *config.Config) *adapters.SubscriptionAdapter {
	subscriptionAdapter, err := adapters.NewSubscriptionAdapter(cfg.SubscriptionServiceGRPCAddr)
	if err != nil {
		panic(err)
	}
	return subscriptionAdapter
}
