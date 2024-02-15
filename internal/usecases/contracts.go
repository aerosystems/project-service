package usecases

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
)

type SubsRepository interface {
	GetSubscription(userUuid uuid.UUID) (models.KindSubscription, int, error)
}

type ProjectRepository interface {
	GetById(Id int) (*models.Project, error)
	GetByToken(token string) (*models.Project, error)
	GetByUserUuid(userUuid uuid.UUID) ([]models.Project, error)
	Create(project *models.Project) error
	Update(project *models.Project) error
	Delete(project *models.Project) error
}
