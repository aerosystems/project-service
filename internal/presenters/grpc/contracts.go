package GRPCServer

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	DetermineStrategy(customerUuid string, role string) error
	CreateDefaultProject(customerUuid uuid.UUID) (*models.Project, error)
	DeleteProject(projectUuid string) error
}
