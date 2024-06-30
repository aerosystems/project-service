package project

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ProjectDelete godoc
// @Summary delete project by ProjectId
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param	projectId	path	string	true "ProjectId"
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects/{projectUuid} [delete]
func (ph Handler) ProjectDelete(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	projectUuid := c.Param("projectUuid")
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return ph.ErrorResponse(c, http.StatusForbidden, "creating project is forbidden", err)
	}
	if project, err := ph.projectUsecase.GetProjectByUuid(projectUuid); err != nil && project == nil {
		return ph.ErrorResponse(c, http.StatusNotFound, "project not found", err)
	} else {
		return ph.ErrorResponse(c, http.StatusForbidden, "user does not have access to this project", err)
	}
	if err := ph.projectUsecase.DeleteProjectByUuid(projectUuid); err != nil {
		return ph.ErrorResponse(c, http.StatusInternalServerError, "could not delete project", err)
	}
	return ph.SuccessResponse(c, http.StatusOK, "project successfully deleted", nil)
}
