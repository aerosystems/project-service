package HTTPServer

import (
	"github.com/aerosystems/project-service/internal/entities"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	*BaseHandler
	projectUsecase ProjectUsecase
}

func NewProjectHandler(baseHandler *BaseHandler, projectUsecase ProjectUsecase) *ProjectHandler {
	return &ProjectHandler{
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

func (p *Project) ToModel() *entities.Project {
	return &entities.Project{
		Uuid:         p.Uuid,
		CustomerUUID: p.CustomerUuid,
		Name:         p.Name,
		Token:        p.Token,
	}
}

func ModelToProject(project *entities.Project) *Project {
	return &Project{
		Uuid:         project.Uuid,
		CustomerUuid: project.CustomerUUID,
		Name:         project.Name,
		Token:        project.Token,
	}
}

func ModelListToProjectList(projects []entities.Project) []Project {
	projectList := make([]Project, 0, len(projects))
	for _, project := range projects {
		projectList = append(projectList, *ModelToProject(&project))
	}
	return projectList
}
