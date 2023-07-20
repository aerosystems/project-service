package models

import (
	"time"
)

type Project struct {
	ID         int       `json:"id" gorm:"primaryKey;unique;autoIncrement" example:"66"`
	UserID     int       `json:"user_id" gorm:"index:idx_user_id_name,unique" example:"666"`
	Name       string    `json:"name" gorm:"index:idx_user_id_name,unique" example:"bla-bla-bla.com"`
	Token      string    `json:"token" example:"38fa45ebb919g5d966122bf9g42a38ceb1e4f6eddf1da70ef00afbdc38197d8f"`
	AccessTime time.Time `json:"access_time" example:"2027-03-03T08:15:00Z"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

type ProjectRepository interface {
	FindByID(ID int) (*Project, error)
	FindByToken(token string) (*Project, error)
	FindByUserID(UserID int) ([]Project, error)
	Create(project *Project) error
	Update(project *Project) error
	Delete(project *Project) error
}
