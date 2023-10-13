package transform

import (
	"github.com/aerosystems/project-service/internal/models"
)

type CreateProjectRequest struct {
	UserID int    `json:"userId" example:"66"`
	Name   string `json:"name" example:"bla-bla-bla.com"`
}

func CreateRequest2Model(reqProject CreateProjectRequest) models.Project {
	return models.Project{
		UserID: reqProject.UserID,
		Name:   reqProject.Name,
	}
}
