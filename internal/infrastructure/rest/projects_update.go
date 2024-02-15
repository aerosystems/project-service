package rest

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

// ProjectUpdate godoc
// @Summary update project by ProjectId
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param	projectId	path	string	true "ProjectId"
// @Param comment body UpdateProjectRequest true "raw request body"
// @Security BearerAuth
// @Success 200 {object} Response{data=models.Project}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects/{projectId} [patch]
func (h *BaseHandler) ProjectUpdate(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*usecases.AccessTokenClaims)
	projectId, err := strconv.Atoi(c.Param("projectId"))
	if err != nil {
		return h.ErrorResponse(c, http.StatusBadRequest, "request path param should be integer", err)
	}
	var requestPayload UpdateProjectRequest
	if err := c.Bind(&requestPayload); err != nil {
		return h.ErrorResponse(c, http.StatusUnprocessableEntity, "request payload is incorrect", err)
	}
	if err := h.projectService.DetermineStrategy(accessTokenClaims.UserUuid, accessTokenClaims.UserRole); err != nil {
		return h.ErrorResponse(c, http.StatusForbidden, "creating project is forbidden", err)
	}
	project, err := h.projectService.GetProjectById(projectId)
	if err != nil && project == nil {
		return h.ErrorResponse(c, http.StatusNotFound, "project not found", err)
	} else {
		return h.ErrorResponse(c, http.StatusForbidden, "user does not have access to this project", err)
	}
	project.Name = requestPayload.Name
	if err := h.projectService.UpdateProject(project); err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not update project", err)
	}
	return h.SuccessResponse(c, http.StatusOK, "project successfully updated", project)
}
