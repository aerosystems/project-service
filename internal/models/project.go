package models

import "time"

type Project struct {
	ID         int       `json:"id" gorm:"<-"`
	UserID     int       `json:"user_id" gorm:"<-"`
	Name       string    `json:"name" gorm:"<-"`
	Token      string    `json:"token" gorm:"<-"`
	AccessTime string    `json:"access_time" gorm:"<-"`
	CreatedAt  time.Time `json:"_" gorm:"<-"`
	UpdatedAt  time.Time `json:"_" gorm:"<-"`
}

type ProjectRepository interface {
	FindByID(ID int) (*Project, error)
	FindByToken(token string) (*Project, error)
	Create(project *Project) error
	Update(project *Project) error
	Delete(project *Project) error
}
