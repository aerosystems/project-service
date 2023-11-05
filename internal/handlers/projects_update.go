package handlers

import (
	"errors"
	"fmt"
	"github.com/aerosystems/project-service/internal/helpers"
	AuthService "github.com/aerosystems/project-service/pkg/auth_service"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type UpdateProjectRequest struct {
	Name       string    `json:"name" example:"bla-bla-bla.com"`
	AccessTime time.Time `json:"accessTime" example:"2027-03-03T08:15:00Z"`
}

// ProjectUpdate godoc
// @Summary update project by Project ID
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param	projectID	path	string	true "Project ID"
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
// @Router /v1/projects/{projectID} [patch]
func (h *BaseHandler) ProjectUpdate(w http.ResponseWriter, r *http.Request) {
	// receive AccessToken Claims from context middleware
	accessTokenClaims, ok := r.Context().Value(helpers.ContextKey("accessTokenClaims")).(*AuthService.AccessTokenClaims)
	if !ok {
		err := errors.New("could not get token claims from Access Token")
		_ = WriteResponse(w, http.StatusUnauthorized, NewErrorPayload(401001, "could not get token claims from Access Token", err))
		return
	}
	projectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422002, "request path param should be integer", err))
		return
	}

	var requestPayload UpdateProjectRequest
	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422001, "request payload is incorrect", err))
		return
	}

	project, err := h.projectRepo.GetById(projectID)
	if errors.Is(err, gorm.ErrRecordNotFound) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		err := fmt.Errorf("project ID %d does not exist", projectID)
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404005, "project ID does not exist", err))
		return
	}
	if err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500101, "could not compare new Project with projects", err))
		return
	}

	if requestPayload.Name != "" {
		project.Name = requestPayload.Name
	}

	// restrict access to project for users with role "startup" or "business"
	if project.UserId != accessTokenClaims.UserID && helpers.Contains([]string{"startup", "business"}, accessTokenClaims.UserRole) {
		err := errors.New("user does not have access to this project")
		_ = WriteResponse(w, http.StatusForbidden, NewErrorPayload(403001, err.Error(), err))
		return
	}

	if err = h.projectRepo.Update(project); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) || strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			_ = WriteResponse(w, http.StatusConflict, NewErrorPayload(409105, "user does not have the same project names", err))
			return
		}
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500104, "could not update project", err))
		return
	}

	_ = WriteResponse(w, http.StatusOK, NewResponsePayload("project successfully updated", project))
	return
}
