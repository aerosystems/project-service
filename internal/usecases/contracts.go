package usecases

import (
	"context"
	"github.com/aerosystems/project-service/internal/entities"
	"github.com/google/uuid"
	"time"
)

type SubscriptionAdapter interface {
	GetSubscription(ctx context.Context, customerUuid uuid.UUID) (entities.SubscriptionType, time.Time, error)
}

type ProjectRepository interface {
	GetByUuid(ctx context.Context, uuid uuid.UUID) (*entities.Project, error)
	GetByToken(ctx context.Context, token string) (*entities.Project, error)
	GetByCustomerUuid(ctx context.Context, customerUuid uuid.UUID) ([]entities.Project, error)
	GetByCustomerUuidAndName(ctx context.Context, userUuid uuid.UUID, name string) (*entities.Project, error)
	Create(ctx context.Context, project *entities.Project) error
	Update(ctx context.Context, project *entities.Project) error
	Delete(ctx context.Context, project *entities.Project) error
}

type CheckmailEventsAdapter interface {
	PublishCreateAccessEvent(token, subscriptionType string, accessTime time.Time) error
}
