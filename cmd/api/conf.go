package main

import (
	"github.com/aerosystems/project-service/internal/handlers"
	"github.com/aerosystems/project-service/internal/models"
)

type Config struct {
	BaseHandler *handlers.BaseHandler
	ProjectRepo models.ProjectRepository
}
