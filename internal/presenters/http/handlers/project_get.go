package handlers

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
// @Success 200 {object} Project
// @Failure 400 {object} echo.HTTPError
// @Failure 401 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 422 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/projects/{projectUuid} [get]
func (ph ProjectHandler) GetProject(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	projectUuid := c.Param("projectUuid")
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Getting project is forbidden.")
	}
	project, err := ph.projectUsecase.GetProjectByUuid(projectUuid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, project)
}
