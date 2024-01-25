package main

import (
	"fmt"
	"github.com/aerosystems/project-service/internal/handlers"
	"github.com/aerosystems/project-service/internal/middleware"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/aerosystems/project-service/internal/repository"
	RPCServer "github.com/aerosystems/project-service/internal/rpc_server"
	RPCServices "github.com/aerosystems/project-service/internal/rpc_services"
	"github.com/aerosystems/project-service/internal/services"
	"github.com/aerosystems/project-service/pkg/gorm_postgres"
	"github.com/aerosystems/project-service/pkg/logger"
	RPCClient "github.com/aerosystems/project-service/pkg/rpc_client"
	"github.com/aerosystems/project-service/pkg/validators"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/rpc"
	"os"
)

const (
	rpcPort = 5001
	webPort = 80
)

// @title Project Service
// @version 1.0.6
// @description A part of microservice infrastructure, who responsible for managing user Projects

// @contact.name Artem Kostenko
// @contact.url https://github.com/aerosystems

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey X-Api-Key
// @in header
// @name X-Api-Key
// @description Should contain Token, digits and letters, 64 symbols length

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Should contain Access JWT Token, with the Bearer started

// @host gw.verifire.app/project
// @schemes https
// @BasePath /
func main() {
	log := logger.NewLogger(os.Getenv("HOSTNAME"))

	clientGORM := GormPostgres.NewClient(logrus.NewEntry(log.Logger))
	_ = clientGORM.AutoMigrate(&models.Project{})

	projectRepo := repository.NewProjectRepo(clientGORM)

	subsRPCClient := RPCClient.NewClient("tcp", "subs-service:5001")
	subsRPC := RPCServices.NewSubsRPC(subsRPCClient)

	projectService := services.NewProjectServiceImpl(projectRepo, subsRPC)

	baseHandler := handlers.NewBaseHandler(os.Getenv("APP_ENV"), log.Logger, projectService)
	projectServer := RPCServer.NewProjectServer(rpcPort, log.Logger, projectService)

	accessTokenService := services.NewAccessTokenServiceImpl(os.Getenv("ACCESS_SECRET"))

	oauthMiddleware := middleware.NewOAuthMiddlewareImpl(accessTokenService)
	basicAuthMiddleware := middleware.NewBasicAuthMiddlewareImpl(os.Getenv("BASIC_AUTH_DOCS_USERNAME"), os.Getenv("BASIC_AUTH_DOCS_PASSWORD"))

	app := NewConfig(baseHandler, oauthMiddleware, basicAuthMiddleware)
	e := app.NewRouter()
	middleware.AddLog(e, log.Logger)

	validator := validator.New()
	e.Validator = &validators.CustomValidator{Validator: validator}

	errChan := make(chan error)

	go func() {
		log.Infof("starting project-service RPC server on port %d\n", rpcPort)
		errChan <- rpc.Register(projectServer)
		errChan <- projectServer.Listen(rpcPort)
	}()

	go func() {
		log.Infof("starting HTTP server project-service on port %d\n", webPort)
		errChan <- e.Start(fmt.Sprintf(":%d", webPort))
	}()

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
