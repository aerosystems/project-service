package models

import (
	"github.com/google/uuid"
	"time"
)

type Project struct {
	Uuid         uuid.UUID
	CustomerUuid uuid.UUID
	Name         string
	Token        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
