package project

import (
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"github.com/aerosystems/project-service/internal/models"
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
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	uuidStr := c.QueryParam("userUuid")
	if len(uuidStr) == 0 {
		uuidStr = accessTokenClaims.UserUuid
	}
	filterUserUuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return CustomErrors.ErrProjectUuidInvalid
	}
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return echo.NewHTTPError(http.StatusForbidden, "Getting projects is forbidden")
	}
	userUuid, err := uuid.Parse(accessTokenClaims.UserUuid)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "user uuid is incorrect")
	}
	projectList, err := ph.projectUsecase.GetProjectListByCustomerUuid(userUuid, filterUserUuid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ModelListToProjectList(projectList))
}
