package handlers

import (
	"github.com/aerosystems/project-service/internal/services"
	"github.com/labstack/echo/v4"
	"net/http"
)

// ProjectCreate godoc
// @Summary create project
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param comment body CreateProjectRequest true "raw request body"
// @Security BearerAuth
// @Success 201 {object} Response{data=models.Project}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects [post]
func (h *BaseHandler) ProjectCreate(c echo.Context) error {
	accessTokenClaims, _ := c.Get("accessTokenClaims").(*services.AccessTokenClaims)
	var requestPayload CreateProjectRequest
	if err := c.Bind(&requestPayload); err != nil {
		return h.ErrorResponse(c, http.StatusUnprocessableEntity, "request payload is incorrect", err)
	}
	if err := h.projectService.DetermineStrategy(accessTokenClaims.UserId, accessTokenClaims.UserRole); err != nil {
		return h.ErrorResponse(c, http.StatusForbidden, "creating project is forbidden", err)
	}
	if err := h.projectService.CreateProject(requestPayload.UserId, requestPayload.Name); err != nil {
		return h.ErrorResponse(c, http.StatusInternalServerError, "could not create default project", err)
	}
	return h.SuccessResponse(c, http.StatusCreated, "project successfully created", nil)
}
