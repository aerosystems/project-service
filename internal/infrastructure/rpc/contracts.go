package RPCServer

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	DetermineStrategy(userUuid string, role string) error
	CreateDefaultProject(userUuid uuid.UUID) error
	GetProjectByToken(token string) (*models.Project, error)
	GetProjectListByUserUuid(userUuid, filterUserUuid uuid.UUID) (projectList []models.Project, err error)
}
