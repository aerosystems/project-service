package handlers

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	InitProject(userUuidStr string) (*models.Project, error)
	DetermineStrategy(userUuidStr string, role string) error
	GetProjectById(projectId int) (*models.Project, error)
	GetProjectListByCustomerUuid(customerUuid, filterUserUuid uuid.UUID) (projectList []models.Project, err error)
	CreateProject(customerUuid uuid.UUID, name string) error
	UpdateProject(project *models.Project) error
	DeleteProjectById(projectId int) error
}

type TokenUsecase interface {
	IsProjectExistByToken(token string) bool
}
