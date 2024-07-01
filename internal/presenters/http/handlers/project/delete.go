package project

import (
	"errors"
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
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
		return ph.ErrorResponse(c, http.StatusForbidden, "Creating a project is forbidden.", err)
	}
	if err := ph.projectUsecase.DeleteProjectByUuid(projectUuid); err != nil {
		var apiErr CustomErrors.ApiError
		if errors.As(err, &apiErr) {
			return ph.ErrorResponse(c, apiErr.HttpCode, apiErr.Message, err)
		}
		return ph.ErrorResponse(c, http.StatusInternalServerError, "Unable to delete the project.", err)
	}
	return ph.SuccessResponse(c, http.StatusOK, "Project has been successfully deleted.", nil)
}
