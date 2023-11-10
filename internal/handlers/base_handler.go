package handlers

import (
	"github.com/aerosystems/project-service/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"strings"
)

type BaseHandler struct {
	mode           string
	log            *logrus.Logger
	projectService services.ProjectService
}

func NewBaseHandler(mode string, log *logrus.Logger, projectService services.ProjectService) *BaseHandler {
	return &BaseHandler{
		mode:           mode,
		log:            log,
		projectService: projectService,
	}
}

type CreateProjectRequest struct {
	UserId int    `json:"userId" validate:"required,number" example:"66"`
	Name   string `json:"name" validate:"required,min=3,max=128" example:"bla-bla-bla.com"`
}

type UpdateProjectRequest struct {
	Name string `json:"name" validate:"required,min=3,max=128" example:"bla-bla-bla.com"`
}

// Response is the type used for sending JSON around
type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ErrorResponse is the type used for sending JSON around
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
}

// SuccessResponse takes a response status code and arbitrary data and writes a json response to the client
func (h *BaseHandler) SuccessResponse(c echo.Context, statusCode int, message string, data any) error {
	payload := Response{
		Message: message,
		Data:    data,
	}
	return c.JSON(statusCode, payload)
}

// ErrorResponse takes a response status code and arbitrary data and writes a json response to the client. It depends on the mode whether the error is included in the response.
func (h *BaseHandler) ErrorResponse(c echo.Context, statusCode int, message string, err error) error {
	payload := Response{Message: message}
	if strings.ToLower(h.mode) == "dev" {
		payload.Data = err.Error()
	}
	return c.JSON(statusCode, payload)
}
