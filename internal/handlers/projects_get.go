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
)

// GetProjectList godoc
// @Summary get all projects. Result depends on user role
// @Tags projects
// @Accept  json
// @Produce application/json
// @Security BearerAuth
// @Success 200 {object} Response{data=[]models.Project}
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects [get]
func (h *BaseHandler) GetProjectList(w http.ResponseWriter, r *http.Request) {
	// receive AccessToken Claims from context middleware
	accessTokenClaims, ok := r.Context().Value(helpers.ContextKey("accessTokenClaims")).(*AuthService.AccessTokenClaims)
	if !ok {
		err := errors.New("could not get token claims from Access Token")
		_ = WriteResponse(w, http.StatusUnauthorized, NewErrorPayload(401001, "could not get token claims from Access Token", err))
		return
	}

	projects, err := h.projectRepo.FindByUserID(accessTokenClaims.UserID)
	if err != nil {
		err := errors.New("could not get projects by UserID")
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500001, "could not get projects by UserID", err))
		return
	}
	if len(projects) == 0 {
		err := errors.New("projects not found")
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404001, "projects not found", err))
		return
	}

	payload := NewResponsePayload("projects successfully found", projects)
	_ = WriteResponse(w, http.StatusOK, payload)
	return
}

// GetProject godoc
// @Summary get project by Project ID
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param	projectID	path	string	true "Project ID"
// @Security BearerAuth
// @Success 200 {object} Response{data=models.Project}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects/{projectID} [get]
func (h *BaseHandler) GetProject(w http.ResponseWriter, r *http.Request) {
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

	project, err := h.projectRepo.FindByID(projectID)
	if err != nil && err != gorm.ErrRecordNotFound {
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(500001, "could not find Project by Project ID", err))
		return
	}

	if project == nil {
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404001, "project not found", err))
		return
	}

	if project.UserID != accessTokenClaims.UserID {
		err := errors.New("project does not belong to user")
		_ = WriteResponse(w, http.StatusForbidden, NewErrorPayload(403001, "project does not belong to user", err))
		return
	}

	_ = WriteResponse(w, http.StatusOK, NewResponsePayload(fmt.Sprintf("project ID %d successfully found", projectID), project))
	return
}
