package handlers

import (
	"github.com/aerosystems/project-service/internal/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// GetProjectList godoc
// @Summary get all projects. Result depends on user role
// @Tags projects
// @Accept  json
// @Produce application/json
// @Security BearerAuth
// @Param	userUuid	query	string	false "UserUuid"
// @Success 200 {object} Response{data=[]models.Project}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects [get]
func (h *BaseHandler) GetProjectList(c echo.Context) (err error) {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*services.AccessTokenClaims)
	uuidStr := c.QueryParam("userUuid")
	if len(uuidStr) == 0 {
		uuidStr = accessTokenClaims.UserUuid
	}
	filterUserUuid, err := uuid.Parse(uuidStr)
	if err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, "user uuid is incorrect", err)
	}
	if err := h.projectService.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return h.ErrorResponse(c, http.StatusForbidden, "creating project is forbidden", err)
	}
	userUuid, err := uuid.Parse(accessTokenClaims.UserUuid)
	if err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, "user uuid is incorrect", err)
	}
	projectList, err := h.projectService.GetProjectListByUserUuid(userUuid, filterUserUuid)
	if err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not get projects", err)
	}
	return h.SuccessResponse(c, http.StatusOK, "projects successfully found", projectList)
}

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
// @Router /v1/projects/{projectId} [get]
func (h *BaseHandler) GetProject(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*services.AccessTokenClaims)
	projectId, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, "request path param should be integer", err)
	}
	if err := h.projectService.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return h.ErrorResponse(c, http.StatusForbidden, "creating project is forbidden", err)
	}
	project, err := h.projectService.GetProjectById(projectId)
	if err != nil {
		return h.ErrorResponse(c, http.StatusForbidden, "user does not have access to this project", err)
	}
	if project == nil {
		return h.ErrorResponse(c, http.StatusNotFound, "project not found", nil)
	}
	return h.SuccessResponse(c, http.StatusOK, "project successfully found", project)
}
