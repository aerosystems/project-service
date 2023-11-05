package handlers

import (
	"errors"
	"fmt"
	"github.com/aerosystems/project-service/internal/helpers"
	"github.com/aerosystems/project-service/internal/transform"
	AuthService "github.com/aerosystems/project-service/pkg/auth_service"
	"github.com/aerosystems/project-service/pkg/validators"
	"net/http"
)

// ProjectCreate godoc
// @Summary create project
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param comment body transform.CreateProjectRequest true "raw request body"
// @Security BearerAuth
// @Success 200 {object} Response{data=models.Project}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /v1/projects [post]
func (h *BaseHandler) ProjectCreate(w http.ResponseWriter, r *http.Request) {
	// receive AccessToken Claims from context middleware
	accessTokenClaims, ok := r.Context().Value(helpers.ContextKey("accessTokenClaims")).(*AuthService.AccessTokenClaims)
	if !ok {
		err := errors.New("could not get token claims from Access Token")
		_ = WriteResponse(w, http.StatusUnauthorized, NewErrorPayload(401001, "could not get token claims from Access Token", err))
		return
	}

	var requestPayload transform.CreateProjectRequest
	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422001, "could not read request body", err))
		return
	}

	newProject := transform.CreateRequest2Model(requestPayload)

	if err := validators.ValidateProject(newProject); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422001, err.Error(), err))
		return
	}

	projectList, err := h.projectRepo.GetByUserId(requestPayload.UserID)
	if err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500101, "could not create new Project", err))
		return
	}

	// restrict User with Role "startup"
	if len(projectList) > 0 && accessTokenClaims.UserRole == "startup" {
		err := fmt.Errorf("user with Startup role and ID %d already has project", requestPayload.UserID)
		_ = WriteResponse(w, http.StatusConflict, NewErrorPayload(409101, "user with Startup plan already has project, for create more projects you should switch into Business plan", err))
		return
	}

	for _, project := range projectList {
		if project.Name == requestPayload.Name {
			err := fmt.Errorf("project with Name %s already exists", requestPayload.Name)
			_ = WriteResponse(w, http.StatusConflict, NewErrorPayload(409102, err.Error(), err))
			return
		}
	}

	if err := h.projectRepo.Create(&newProject); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500103, "could not create new Project", err))
		return
	}

	_ = WriteResponse(w, http.StatusOK, NewResponsePayload("project successfully created", newProject))
	return
}
