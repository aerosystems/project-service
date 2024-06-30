package project

import (
	"github.com/aerosystems/project-service/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetProject godoc
// @Summary get project by ProjectId
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param	projectId	path	string	true "ProjectId"
// @Security BearerAuth
// @Success 200 {object} Response{data=models.Project}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects/{projectUuid} [get]
func (ph Handler) GetProject(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	projectUuid := c.Param("projectUuid")
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return ph.ErrorResponse(c, http.StatusForbidden, "getting project is forbidden", err)
	}
	project, err := ph.projectUsecase.GetProjectByUuid(projectUuid)
	if err != nil {
		return ph.ErrorResponse(c, http.StatusForbidden, "user does not have access to this project", err)
	}
	if project == nil {
		return ph.ErrorResponse(c, http.StatusNotFound, "project not found", nil)
	}
	return ph.SuccessResponse(c, http.StatusOK, "project successfully found", project)
}
