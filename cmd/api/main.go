package main

import (
	"fmt"
	"github.com/aerosystems/project-service/internal/handlers"
	"github.com/aerosystems/project-service/internal/repository"
	"github.com/aerosystems/project-service/pkg/gorm_postgres"
	"log"
	"net/http"
)

// @title Project Service
// @version 1.0
// @description A part of microservice infrastructure, who responsible for managing user Projects

// @contact.name Artem Kostenko
// @contact.url https://github.com/aerosystems

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @securityDefinitions.apikey X-API-KEY
// @in header
// @name X-API-KEY
// @description Should contain Token, digits and letters, 64 symbols length

// @host localhost:8082
// @BasePath /
func main() {
	clientGORM := GormPostgres.NewClient()
	projectRepo := repository.NewProjectRepo(clientGORM)

	app := Config{
		WebPort:     "80",
		BaseHandler: handlers.NewBaseHandler(projectRepo),
		ProjectRepo: projectRepo,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", app.WebPort),
		Handler: app.routes(),
	}

	log.Printf("Starting authentication end service on port %s\n", app.WebPort)
	err := srv.ListenAndServe()

	if err != nil {
		log.Panic(err)
	}
}
