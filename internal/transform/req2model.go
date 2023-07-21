package transform

import (
	"github.com/aerosystems/project-service/internal/models"
	"time"
)

type CreateProjectRequest struct {
	UserID     int       `json:"userId" example:"66"`
	Name       string    `json:"name" example:"bla-bla-bla.com"`
	AccessTime time.Time `json:"accessTime" example:"2027-03-03T08:15:00Z"`
}

func CreateRequest2Model(reqProject CreateProjectRequest) models.Project {
	return models.Project{
		UserID:     reqProject.UserID,
		Name:       reqProject.Name,
		AccessTime: reqProject.AccessTime,
	}
}
