package HTTPServer

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"time"
)

type ProjectUsecase interface {
	InitProject(userUuidStr, subscriptionType string, accessTime time.Time) (*models.Project, error)
	DetermineStrategy(userUuidStr string, role string) error
	GetProjectByUuid(projectUuidStr string) (*models.Project, error)
	GetProjectListByCustomerUuid(customerUuid, filterUserUuid uuid.UUID) (projectList []models.Project, err error)
	CreateProject(customerUuid uuid.UUID, name string) (*models.Project, error)
	UpdateProject(projectUuidStr, projectName string) (*models.Project, error)
	DeleteProject(projectUuidStr string) error
}

type TokenUsecase interface {
	IsProjectExistByToken(token string) bool
}
