package HTTPServer

import (
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/aerosystems/project-service/internal/models"
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

// GetProjectList godoc
// @Summary get all projects. Result depends on user role
// @Tags projects
// @Accept  json
// @Produce application/json
// @Security BearerAuth
// @Param	userUuid	query	string	false "CustomerUUID"
// @Success 200 {object} []Project
// @Failure 401 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/projects [get]
func (ph ProjectHandler) GetProjectList(c echo.Context) (err error) {
	userClaims, err := GetUserClaimsFromContext(c.Request().Context())
	if err != nil {
		return err
	}
	filterUserUuid := userClaims.Uuid
	if len(c.QueryParam("userUuid")) != 0 {
		filterUserUuid, err = uuid.Parse(c.QueryParam("userUuid"))
		if err != nil {
			return CustomErrors.ErrProjectUuidInvalid
		}
	}
	if err := ph.projectUsecase.DetermineStrategy(userClaims.Uuid.String(), userClaims.Role.String()); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Getting projects is forbidden")
	}
	projectList, err := ph.projectUsecase.GetProjectListByCustomerUuid(userClaims.Uuid, filterUserUuid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ModelListToProjectList(projectList))
}
