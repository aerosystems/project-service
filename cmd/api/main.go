package main

import (
	"fmt"
	"github.com/aerosystems/project-service/internal/handlers"
	"github.com/aerosystems/project-service/internal/middleware"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/aerosystems/project-service/internal/repository"
	RPCServer "github.com/aerosystems/project-service/internal/rpc_server"
	"github.com/aerosystems/project-service/internal/services"
	"github.com/aerosystems/project-service/pkg/gorm_postgres"
	"github.com/aerosystems/project-service/pkg/logger"
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

// @host gw.verifire.com/project
// @schemes https
// @BasePath /
func main() {
	log := logger.NewLogger(os.Getenv("HOSTNAME"))

	clientGORM := GormPostgres.NewClient(logrus.NewEntry(log.Logger))
	_ = clientGORM.AutoMigrate(&models.Project{})

	projectRepo := repository.NewProjectRepo(clientGORM)
	projectService := services.NewProjectServiceImpl(projectRepo)

	baseHandler := handlers.NewBaseHandler(os.Getenv("APP_ENV"), log.Logger, projectService)
	projectServer := RPCServer.NewProjectServer(rpcPort, log.Logger, projectService)

	app := NewConfig(baseHandler)
	e := app.NewRouter()
	middleware.AddMiddleware(e, log.Logger)
	validator := validator.New()
	e.Validator = &validators.CustomValidator{Validator: validator}

	errChan := make(chan error)

	go func() {
		log.Infof("starting project-service RPC server on port %d\n", rpcPort)
		errChan <- rpc.Register(projectServer)
		errChan <- projectServer.Listen(rpcPort)
	}()

	go func() {
		log.Infof("starting HTTP server project-service on port %s\n", webPort)
		errChan <- e.Start(fmt.Sprintf(":%d", webPort))
	}()

	err := <-errChan
	log.Fatal(err)
}
