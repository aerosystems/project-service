package project

import (
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/aerosystems/project-service/internal/models"
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
// @Success 201 {object} Response{data=models.Project}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects [post]
func (ph Handler) ProjectCreate(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	var requestPayload CreateProjectRequest
	if err := c.Bind(&requestPayload); err != nil {
		return ph.ErrorResponse(c, CustomErrors.ErrRequestPayloadIncorrect.HttpCode, CustomErrors.ErrRequestPayloadIncorrect.Message, err)
	}
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return ph.ErrorResponse(c, http.StatusForbidden, "creating project is forbidden", err)
	}
	userUuid, err := uuid.Parse(requestPayload.UserUuid)
	if err != nil {
		return ph.ErrorResponse(c, http.StatusBadRequest, "user uuid is incorrect", err)
	}
	if err := ph.projectUsecase.CreateProject(userUuid, requestPayload.Name); err != nil {
		return ph.ErrorResponse(c, http.StatusInternalServerError, "could not create default project", err)
	}
	return ph.SuccessResponse(c, http.StatusCreated, "project successfully created", nil)
}
