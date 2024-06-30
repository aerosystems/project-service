package project

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/aerosystems/project-service/internal/presenters/http/handlers"
	"github.com/google/uuid"
)

type Handler struct {
	*handlers.BaseHandler
	projectUsecase handlers.ProjectUsecase
}

func NewProjectHandler(baseHandler *handlers.BaseHandler, projectUsecase handlers.ProjectUsecase) *Handler {
	return &Handler{
		BaseHandler:    baseHandler,
		projectUsecase: projectUsecase,
	}
}

type Project struct {
	Uuid         uuid.UUID `json:"uuid" example:"8893ef16-0030-4f90-a686-b82a6a57ef93"`
	CustomerUuid uuid.UUID `json:"customerUuid" example:"666"`
	Name         string    `json:"name" example:"bla-bla-bla.com"`
	Token        string    `json:"token" example:"38fa45ebb919g5d966122bf9g42a38ceb1e4f6eddf1da70ef00afbdc38197d8f"`
}

func (p *Project) ToModel() *models.Project {
	return &models.Project{
		Uuid:         p.Uuid,
		CustomerUuid: p.CustomerUuid,
		Name:         p.Name,
		Token:        p.Token,
	}
}

func ModelToProject(project *models.Project) *Project {
	return &Project{
		Uuid:         project.Uuid,
		CustomerUuid: project.CustomerUuid,
		Name:         project.Name,
		Token:        project.Token,
	}
}
