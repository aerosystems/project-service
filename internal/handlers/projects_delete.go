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

// ProjectDelete godoc
// @Summary delete project by Project ID
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param	projectID	path	string	true "Project ID"
// @Security BearerAuth
// @Success 200 {object} Response
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects/{projectID} [delete]
func (h *BaseHandler) ProjectDelete(w http.ResponseWriter, r *http.Request) {
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

	if err == gorm.ErrRecordNotFound {
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404001, "project not found", err))
		return
	}
	if err != nil {
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(500001, "could not find Project by ProjectID", err))
		return
	}

	// restrict access to project for users with role "startup" or "business"
	if project.UserID != accessTokenClaims.UserID && helpers.Contains([]string{"startup", "business"}, accessTokenClaims.UserRole) {
		err := errors.New("user does not have access to this project")
		_ = WriteResponse(w, http.StatusForbidden, NewErrorPayload(403001, err.Error(), err))
		return
	}

	if err = h.projectRepo.Delete(project); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500106, "could not delete Project", err))
		return
	}

	_ = WriteResponse(w, http.StatusOK, NewResponsePayload(fmt.Sprintf("project ID %d successfuly deleted", projectID), nil))
	return
}
