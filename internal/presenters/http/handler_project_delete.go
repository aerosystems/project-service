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
// @Success 204
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects/{projectUuid} [delete]
func (ph ProjectHandler) ProjectDelete(c echo.Context) error {
	user, err := GetUserFromContext(c.Request().Context())
	if err != nil {
		return err
	}

	if err = ph.projectUsecase.DetermineStrategy(c.Request().Context(), user.UUID, user.Role); err != nil {
		return err
	}

	projectUUID := c.Param("projectUuid")
	if err = ph.projectUsecase.DeleteProject(c.Request().Context(), projectUUID); err != nil {
		return err
	}
	return c.JSON(http.StatusNoContent, nil)
}
