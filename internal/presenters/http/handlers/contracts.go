package handlers

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	InitProject(userUuidStr string) (*models.Project, error)
	DetermineStrategy(userUuidStr string, role string) error
	GetProjectByUuid(projectUuidStr string) (*models.Project, error)
	GetProjectListByCustomerUuid(customerUuid, filterUserUuid uuid.UUID) (projectList []models.Project, err error)
	CreateProject(customerUuid uuid.UUID, name string) error
	UpdateProject(projectUuidStr, projectName string) (*models.Project, error)
	DeleteProjectByUuid(projectUuidStr string) error
}

type TokenUsecase interface {
	IsProjectExistByToken(token string) bool
}
