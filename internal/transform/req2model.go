package transform

import (
	"github.com/aerosystems/project-service/internal/models"
)

type CreateProjectRequest struct {
	UserId int    `json:"userId" example:"66"`
	Name   string `json:"name" example:"bla-bla-bla.com"`
}

func CreateRequest2Model(reqProject CreateProjectRequest) models.Project {
	return models.Project{
		UserId: reqProject.UserId,
		Name:   reqProject.Name,
	}
}
