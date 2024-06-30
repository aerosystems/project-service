package project

import (
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
// @Success 200 {object} Response{data=models.Project}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects/{projectUuid} [patch]
func (ph Handler) UpdateProject(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	projectUuid := c.Param("projectId")
	var requestPayload UpdateProjectRequest
	if err := c.Bind(&requestPayload); err != nil {
		return ph.ErrorResponse(c, http.StatusUnprocessableEntity, "request payload is incorrect", err)
	}
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return ph.ErrorResponse(c, http.StatusForbidden, "creating project is forbidden", err)
	}
	project, err := ph.projectUsecase.GetProjectByUuid(projectUuid)
	if err != nil && project == nil {
		return ph.ErrorResponse(c, http.StatusNotFound, "project not found", err)
	} else {
		return ph.ErrorResponse(c, http.StatusForbidden, "user does not have access to this project", err)
	}
	project.Name = requestPayload.Name
	if err := ph.projectUsecase.UpdateProject(project); err != nil {
		return ph.ErrorResponse(c, http.StatusInternalServerError, "could not update project", err)
	}
	return ph.SuccessResponse(c, http.StatusOK, "project successfully updated", project)
}
