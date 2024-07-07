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
// @Success 201 {object} Project
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 422 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/projects [post]
func (ph Handler) ProjectCreate(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	var requestPayload CreateProjectRequest
	if err := c.Bind(&requestPayload); err != nil {
		return CustomErrors.ErrInvalidRequestBody
	}
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "creating project is forbidden")
	}
	userUuid, err := uuid.Parse(requestPayload.UserUuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user uuid is incorrect")
	}
	project, err := ph.projectUsecase.CreateProject(userUuid, requestPayload.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, ModelToProject(project))
}
