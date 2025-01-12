package usecases

import (
	"context"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/google/uuid"
	"time"
)

type SubscriptionAdapter interface {
	GetSubscription(ctx context.Context, customerUuid uuid.UUID) (models.SubscriptionType, time.Time, error)
}

type ProjectRepository interface {
	GetByUuid(ctx context.Context, uuid uuid.UUID) (*models.Project, error)
	GetByToken(ctx context.Context, token string) (*models.Project, error)
	GetByCustomerUuid(ctx context.Context, customerUuid uuid.UUID) ([]models.Project, error)
	GetByCustomerUuidAndName(ctx context.Context, userUuid uuid.UUID, name string) (*models.Project, error)
	Create(ctx context.Context, project *models.Project) error
	Update(ctx context.Context, project *models.Project) error
	Delete(ctx context.Context, project *models.Project) error
}

type CheckmailEventsAdapter interface {
	PublishCreateAccessEvent(token, subscriptionType string, accessTime time.Time) error
}
