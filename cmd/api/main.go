package main

import (
	"fmt"
	"github.com/aerosystems/project-service/internal/handlers"
	"github.com/aerosystems/project-service/internal/repository"
	"github.com/aerosystems/project-service/pkg/mygorm"
	"log"
	"net/http"
)

const webPort = "80"

type Config struct {
	BaseHandler *handlers.BaseHandler
}

// @title Project Service
// @version 1.0
// @description A part of microservice infrastructure, who responsible for managing user Projects

// @contact.name Artem Kostenko
// @contact.url https://github.com/aerosystems

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
// @BasePath /v1
func main() {
	clientGORM := mygorm.NewClient()
	projectRepo := repository.NewProjectRepo(clientGORM)

	app := Config{
		BaseHandler: handlers.NewBaseHandler(projectRepo),
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	log.Printf("Starting authentication end service on port %s\n", webPort)
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
