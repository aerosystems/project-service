package models

import (
	"time"
)

type Project struct {
	ID         int       `json:"id" gorm:"<-"`
	UserID     int       `json:"user_id" gorm:"<-"`
	Name       string    `json:"name" gorm:"<-"`
	Token      string    `json:"-" gorm:"<-"`
	AccessTime time.Time `json:"access_time" gorm:"<-"`
	CreatedAt  time.Time `json:"-" gorm:"<-"`
	UpdatedAt  time.Time `json:"-" gorm:"<-"`
}

type ProjectRequest struct {
	UserID     int       `json:"user_id"`
	Name       string    `json:"name"`
	AccessTime time.Time `json:"access_time"`
}

type ProjectRepository interface {
	FindByID(ID int) (*Project, error)
	FindByToken(token string) (*Project, error)
	FindByUserID(UserID int) (*Project, error)
	Create(project *Project) error
	Update(project *Project) error
	Delete(project *Project) error
}
