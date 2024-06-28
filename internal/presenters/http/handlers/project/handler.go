package project

import (
	"github.com/aerosystems/project-service/internal/presenters/http/handlers"
)

type Handler struct {
	*handlers.BaseHandler
	projectUsecase handlers.ProjectUsecase
}

func NewProjectHandler(baseHandler *handlers.BaseHandler, projectUsecase handlers.ProjectUsecase) *Handler {
	return &Handler{
		BaseHandler:    baseHandler,
		projectUsecase: projectUsecase,
	}
}
