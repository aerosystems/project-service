package main

import (
	"fmt"
	"github.com/aerosystems/project-service/internal/handlers"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/aerosystems/project-service/internal/repository"
	RPCServer "github.com/aerosystems/project-service/internal/rpc_server"
	"github.com/aerosystems/project-service/pkg/gorm_postgres"
	"log"
	"net/http"
	"net/rpc"
)

const (
	rpcPort = "5001"
	webPort = "80"
)

// @title Project Service
// @version 1.0.5
// @description A part of microservice infrastructure, who responsible for managing user Projects

// @contact.name Artem Kostenko
// @contact.url https://github.com/aerosystems

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey X-API-KEY
// @in header
// @name X-API-KEY
// @description Should contain Token, digits and letters, 64 symbols length

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Should contain Access JWT Token, with the Bearer started

// @host localhost:8082
// @BasePath /
func main() {
	clientGORM := GormPostgres.NewClient()
	if err := clientGORM.AutoMigrate(&models.Project{}); err != nil {
		log.Panic(err)
	}
	projectRepo := repository.NewProjectRepo(clientGORM)

	app := Config{
		BaseHandler: handlers.NewBaseHandler(projectRepo),
		ProjectRepo: projectRepo,
	}

	if err := rpc.Register(RPCServer.NewProjectServer(projectRepo)); err != nil {
		log.Fatal(err)
	}
	errChan := make(chan error)
	// Start RPC server
	log.Printf("starting RPC server project-service on port %s\n", rpcPort)
	go func() {
		if err := RPCServer.Listen(rpcPort); err != nil {
			errChan <- err
		}
	}()

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}
	// Start HTTP server
	log.Printf("starting HTTP server project-service on port %s\n", webPort)
	go func() {
		errChan <- srv.ListenAndServe()
	}()

	err := <-errChan
	log.Fatal(err)
}
