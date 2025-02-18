package HTTPServer

import (
	"github.com/aerosystems/project-service/internal/entities"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UpdateProjectRequest struct {
	Name string `json:"name" validate:"required,min=3,max=128" example:"bla-bla-bla.com"`
}

// UpdateProject godoc
// @Summary update project by ProjectId
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param	projectId	path	string	true "ProjectId"
// @Param comment body UpdateProjectRequest true "raw request body"
// @Security BearerAuth
// @Success 200 {object} Project
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects/{projectUuid} [patch]
func (h Handler) UpdateProject(c echo.Context) error {
	user, err := GetUserFromContext(c.Request().Context())
	if err != nil {
		return err
	}

	var requestPayload UpdateProjectRequest
	if err = c.Bind(&requestPayload); err != nil {
		return entities.ErrInvalidRequestBody
	}
	if err = h.projectUsecase.DetermineStrategy(c.Request().Context(), user.UUID, user.Role); err != nil {
		return err
	}

	projectUUID := c.Param("project_uuid")
	project, err := h.projectUsecase.UpdateProject(c.Request().Context(), projectUUID, requestPayload.Name)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ModelToProject(project))
}
