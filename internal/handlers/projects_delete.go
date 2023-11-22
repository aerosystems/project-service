package handlers

import (
	"github.com/aerosystems/project-service/internal/services"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
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
// @Router /v1/projects/{projectId} [delete]
func (h *BaseHandler) ProjectDelete(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*services.AccessTokenClaims)
	projectId, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, "request path param should be integer", err)
	}
	if err := h.projectService.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return h.ErrorResponse(c, http.StatusForbidden, "creating project is forbidden", err)
	}
	if project, err := h.projectService.GetProjectById(projectId); err != nil && project == nil {
		return h.ErrorResponse(c, http.StatusNotFound, "project not found", err)
	} else {
		return h.ErrorResponse(c, http.StatusForbidden, "user does not have access to this project", err)
	}
	if err := h.projectService.DeleteProjectById(projectId); err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not delete project", err)
	}
	return h.SuccessResponse(c, http.StatusOK, "project successfully deleted", nil)
}
