package handlers

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	DetermineStrategy(userUuidStr string, role string) error
	GetProjectById(projectId int) (*models.Project, error)
	GetProjectListByUserUuid(userUuid, filterUserUuid uuid.UUID) (projectList []models.Project, err error)
	CreateProject(userUuid uuid.UUID, name string) error
	UpdateProject(project *models.Project) error
	DeleteProjectById(projectId int) error
}

type TokenUsecase interface {
	IsProjectExistByToken(token string) bool
}
