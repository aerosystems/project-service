package models

import (
	"github.com/google/uuid"
	"time"
)

type Project struct {
	Id        int       `json:"id" example:"66"`
	UserUuid  uuid.UUID `json:"userUuid" example:"666"`
	Name      string    `json:"name" example:"bla-bla-bla.com"`
	Token     string    `json:"token" example:"38fa45ebb919g5d966122bf9g42a38ceb1e4f6eddf1da70ef00afbdc38197d8f"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
