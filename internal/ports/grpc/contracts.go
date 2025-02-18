package GRPCServer

import (
	"context"
	"github.com/aerosystems/project-service/internal/entities"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	CreateDefaultProject(ctx context.Context, customerUUID uuid.UUID) (*entities.Project, error)
	DeleteProject(ctx context.Context, projectUUID string) error
}
