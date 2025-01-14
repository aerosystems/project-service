package HTTPServer

import (
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/google/uuid"
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
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects/{projectUuid} [get]
func (ph ProjectHandler) GetProject(c echo.Context) error {
	user, err := GetUserFromContext(c.Request().Context())
	if err != nil {
		return err
	}

	if err = ph.projectUsecase.DetermineStrategy(c.Request().Context(), user.UUID, user.Role); err != nil {
		return err
	}

	projectUUID := c.Param("projectId")
	project, err := ph.projectUsecase.GetProjectByUuid(c.Request().Context(), projectUUID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, project)
}

// GetProjectList godoc
// @Summary get all projects. Result depends on user role
// @Tags projects
// @Accept  json
// @Produce application/json
// @Security BearerAuth
// @Param	userUuid	query	string	false "CustomerUUID"
// @Success 200 {object} []Project
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects [get]
func (ph ProjectHandler) GetProjectList(c echo.Context) (err error) {
	user, err := GetUserFromContext(c.Request().Context())
	if err != nil {
		return err
	}
	var filteredUserUUID uuid.UUID
	if len(c.QueryParam("userUuid")) != 0 {
		filteredUserUUID, err = uuid.Parse(c.QueryParam("userUuid"))
		if err != nil {
			return CustomErrors.ErrProjectUuidInvalid
		}
	}

	if err = ph.projectUsecase.DetermineStrategy(c.Request().Context(), user.UUID, user.Role); err != nil {
		return err
	}

	projectList, err := ph.projectUsecase.GetProjectListByCustomerUuid(c.Request().Context(), user.UUID, filteredUserUUID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ModelListToProjectList(projectList))
}
