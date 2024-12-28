package project

import (
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/aerosystems/project-service/internal/presenters/http/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

// GetProjectList godoc
// @Summary get all projects. Result depends on user role
// @Tags projects
// @Accept  json
// @Produce application/json
// @Security BearerAuth
// @Param	userUuid	query	string	false "CustomerUuid"
// @Success 200 {object} []Project
// @Failure 401 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /v1/projects [get]
func (ph Handler) GetProjectList(c echo.Context) (err error) {
	userClaims, err := middleware.GetUserClaimsFromContext(c.Request().Context())
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
