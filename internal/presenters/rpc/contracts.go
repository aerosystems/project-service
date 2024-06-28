package RpcServer

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	DetermineStrategy(customerUuid string, role string) error
	CreateDefaultProject(customerUuid uuid.UUID) error
	GetProjectByToken(token string) (*models.Project, error)
	GetProjectListByCustomerUuid(customerUuid, filterUserUuid uuid.UUID) (projectList []models.Project, err error)
}
