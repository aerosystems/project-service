package HTTPServer

import (
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
// @Success 204 {object} struct{} "No Content"
// @Failure 401 {object} echo.HTTPError
// @Failure 403 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/projects/{projectUuid} [delete]
func (ph ProjectHandler) ProjectDelete(c echo.Context) error {
	user, err := GetUserFromContext(c.Request().Context())
	if err != nil {
		return err
	}
	projectUuid := c.Param("projectUuid")
	if err := ph.projectUsecase.DetermineStrategy(user.UUID.String(), user.Role.String()); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Deleting a project is forbidden.")
	}
	if err := ph.projectUsecase.DeleteProject(projectUuid); err != nil {
		return err
	}
	return c.JSON(http.StatusNoContent, nil)
}
