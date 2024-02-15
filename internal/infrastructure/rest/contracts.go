package rest

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	DetermineStrategy(userUuidStr string, role string) error
	GetProjectById(projectId int) (*models.Project, error)
	GetProjectByToken(token string) (*models.Project, error)
	GetProjectListByUserUuid(userUuid, filterUserUuid uuid.UUID) (projectList []models.Project, err error)
	CreateProject(userUuid uuid.UUID, name string) error
	CreateDefaultProject(userUuid uuid.UUID) error
	UpdateProject(project *models.Project) error
	DeleteProjectById(projectId int) error
	IsProjectExistByToken(projectToken string) bool
	isProjectNameExist(name string, projectList []models.Project) bool
}
