package models

import (
	"github.com/google/uuid"
	"time"
)

type Project struct {
	Id        int       `json:"id" gorm:"primaryKey;unique;autoIncrement" example:"66"`
	UserUuid  uuid.UUID `json:"userId" gorm:"index:idx_user_id_name,unique" example:"666"`
	Name      string    `json:"name" gorm:"index:idx_user_id_name,unique" example:"bla-bla-bla.com"`
	Token     string    `json:"token" example:"38fa45ebb919g5d966122bf9g42a38ceb1e4f6eddf1da70ef00afbdc38197d8f"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type ProjectRepository interface {
	GetById(Id int) (*Project, error)
	GetByToken(token string) (*Project, error)
	GetByUserUuid(userUuid uuid.UUID) ([]Project, error)
	Create(project *Project) error
	Update(project *Project) error
	Delete(project *Project) error
}
