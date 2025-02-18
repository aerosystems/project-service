package HTTPServer

import (
	"context"
	"github.com/aerosystems/project-service/internal/entities"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	DetermineStrategy(ctx context.Context, userUUID uuid.UUID, role entities.Role) error
	GetProjectByUuid(ctx context.Context, projectUuidStr string) (*entities.Project, error)
	GetProjectListByCustomerUuid(ctx context.Context, customerUuid, filterUserUuid uuid.UUID) (projectList []entities.Project, err error)
	CreateProject(ctx context.Context, customerUuid uuid.UUID, name string) (*entities.Project, error)
	UpdateProject(ctx context.Context, projectUuidStr, projectName string) (*entities.Project, error)
	DeleteProject(ctx context.Context, projectUuidStr string) error
}

type TokenUsecase interface {
	IsProjectExistByToken(token string) bool
}
