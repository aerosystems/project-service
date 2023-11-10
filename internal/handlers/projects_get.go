package handlers

import (
	"github.com/aerosystems/project-service/internal/services"
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
// @Param	userId	query	int	false "UserId"
// @Success 200 {object} Response{data=[]models.Project}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects [get]
func (h *BaseHandler) GetProjectList(c echo.Context) (err error) {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*services.AccessTokenClaims)
	filterUserId, err := strconv.Atoi(c.QueryParam("userId"))
	if err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, "request query param should be integer", err)
	}
	if err := h.projectService.DetermineStrategy(accessTokenClaims.UserId, accessTokenClaims.UserRole); err != nil {
		return h.ErrorResponse(c, http.StatusForbidden, "creating project is forbidden", err)
	}
	projectList, err := h.projectService.GetProjectListByUserId(accessTokenClaims.UserId, filterUserId)
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
	if err := h.projectService.DetermineStrategy(accessTokenClaims.UserId, accessTokenClaims.UserRole); err != nil {
		return h.ErrorResponse(c, http.StatusForbidden, "creating project is forbidden", err)
	}
	project, err := h.projectService.GetProjectById(projectId)
	if err != nil && project == nil {
		return h.ErrorResponse(c, http.StatusNotFound, "project not found", err)
	} else {
		return h.ErrorResponse(c, http.StatusForbidden, "user does not have access to this project", err)
	}
	return h.SuccessResponse(c, http.StatusOK, "project successfully found", project)
}
