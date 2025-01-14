package GRPCServer

import (
	"context"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	CreateDefaultProject(ctx context.Context, customerUUID uuid.UUID) (*models.Project, error)
	DeleteProject(ctx context.Context, projectUUID string) error
}
