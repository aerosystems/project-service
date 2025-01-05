package handlers

import (
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UpdateProjectRequest struct {
	Name string `json:"name" validate:"required,min=3,max=128" example:"bla-bla-bla.com"`
}

// UpdateProject godoc
// @Summary update project by ProjectId
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param	projectId	path	string	true "ProjectId"
// @Param comment body UpdateProjectRequest true "raw request body"
// @Security BearerAuth
// @Success 200 {object} Project
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 409 {object} echo.HTTPError
// @Failure 422 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/projects/{projectUuid} [patch]
func (ph ProjectHandler) UpdateProject(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	projectUuid := c.Param("projectId")
	var requestPayload UpdateProjectRequest
	if err := c.Bind(&requestPayload); err != nil {
		return CustomErrors.ErrInvalidRequestBody
	}
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Updating a project is forbidden.")
	}
	project, err := ph.projectUsecase.UpdateProject(projectUuid, requestPayload.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ModelToProject(project))
}
