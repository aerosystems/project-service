package project

import (
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
// @Success 200 {object} Response{data=[]models.Project}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects [get]
func (ph Handler) GetProjectList(c echo.Context) (err error) {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*models.AccessTokenClaims)
	uuidStr := c.QueryParam("userUuid")
	if len(uuidStr) == 0 {
		uuidStr = accessTokenClaims.UserUuid
	}
	filterUserUuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return ph.ErrorResponse(c, http.StatusBadRequest, "user uuid is incorrect", err)
	}
	if err := ph.projectUsecase.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return ph.ErrorResponse(c, http.StatusForbidden, "getting project is forbidden", err)
	}
	userUuid, err := uuid.Parse(accessTokenClaims.UserUuid)
	if err != nil {
		return ph.ErrorResponse(c, http.StatusBadRequest, "user uuid is incorrect", err)
	}
	projectList, err := ph.projectUsecase.GetProjectListByCustomerUuid(userUuid, filterUserUuid)
	if err != nil {
		return ph.ErrorResponse(c, http.StatusInternalServerError, "could not get projects", err)
	}
	return ph.SuccessResponse(c, http.StatusOK, "projects successfully found", projectList)
}
