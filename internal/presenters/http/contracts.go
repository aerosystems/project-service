package HTTPServer

import (
	"context"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
)

type ProjectUsecase interface {
	DetermineStrategy(ctx context.Context, userUUID uuid.UUID, role models.Role) error
	GetProjectByUuid(ctx context.Context, projectUuidStr string) (*models.Project, error)
	GetProjectListByCustomerUuid(ctx context.Context, customerUuid, filterUserUuid uuid.UUID) (projectList []models.Project, err error)
	CreateProject(ctx context.Context, customerUuid uuid.UUID, name string) (*models.Project, error)
	UpdateProject(ctx context.Context, projectUuidStr, projectName string) (*models.Project, error)
	DeleteProject(ctx context.Context, projectUuidStr string) error
}

type TokenUsecase interface {
	IsProjectExistByToken(token string) bool
}
