package handlers

import (
	"errors"
	"fmt"
	"github.com/aerosystems/project-service/internal/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

// ProjectCreate godoc
// @Summary create project
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param comment body models.ProjectRequest true "raw request body"
// @Param Authorization header string true "should contain Access Token, with the Bearer started"
// @Success 200 {object} Response{data=models.Project}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /projects [post]
func (h *BaseHandler) ProjectCreate(w http.ResponseWriter, r *http.Request) {
	var requestPayload models.ProjectRequest

	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400001, "request payload is incorrect", err))
		return
	}

	if requestPayload.UserID == 0 {
		err := errors.New("claim UserID does not be empty")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400103, err.Error(), err))
		return
	}

	if requestPayload.Name == "" {
		err := errors.New("claim Name does not be empty")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400104, err.Error(), err))
		return
	}

	if requestPayload.AccessTime.IsZero() {
		err := errors.New("claim Access Time does not be empty")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400104, err.Error(), err))
		return
	}

	if requestPayload.AccessTime.Before(time.Now()) {
		err := errors.New("claim Access Time should be more then NOW")
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400102, err.Error(), err))
		return
	}

	project, err := h.projectRepo.FindByUserID(requestPayload.UserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500101, "could not compare new Project with projects in storage", err))
		return
	}

	if project != nil {
		if project.Name == requestPayload.Name {
			err := fmt.Errorf("project with claim Name %s already exists", requestPayload.Name)
			_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500102, "project with claim Name already exists", err))
			return
		}
	}

	var newProject = models.Project{
		UserID:     requestPayload.UserID,
		Name:       requestPayload.Name,
		AccessTime: requestPayload.AccessTime,
	}

	if err = h.projectRepo.Create(&newProject); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500103, "could not create new Project", err))
		return
	}

	payload := NewResponsePayload("project successfully created", newProject)
	_ = WriteResponse(w, http.StatusOK, payload)
	return
}
