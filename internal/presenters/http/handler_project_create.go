package HTTPServer

import (
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type CreateProjectRequest struct {
	UserUuid string `json:"userUuid" validate:"required,number" example:"66"`
	Name     string `json:"name" validate:"required,min=3,max=128" example:"bla-bla-bla.com"`
}

// ProjectCreate godoc
// @Summary create project
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param comment body CreateProjectRequest true "raw request body"
// @Security BearerAuth
// @Success 201 {object} Project
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects [post]
func (ph ProjectHandler) ProjectCreate(c echo.Context) error {
	user, err := GetUserFromContext(c.Request().Context())
	if err != nil {
		return err
	}
	var requestPayload CreateProjectRequest
	if err = c.Bind(&requestPayload); err != nil {
		return CustomErrors.ErrInvalidRequestBody
	}
	if err = ph.projectUsecase.DetermineStrategy(c.Request().Context(), user.UUID, user.Role); err != nil {
		return err
	}
	userUuid, err := uuid.Parse(requestPayload.UserUuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user uuid is incorrect")
	}
	project, err := ph.projectUsecase.CreateProject(c.Request().Context(), userUuid, requestPayload.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, ModelToProject(project))
}
