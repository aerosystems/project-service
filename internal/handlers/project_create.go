package handlers

import (
	"errors"
	"fmt"
	"github.com/aerosystems/project-service/internal/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type CreateProjectRequest struct {
	UserID     int       `json:"user_id" example:"66"`
	Name       string    `json:"name" example:"bla-bla-bla.com"`
	AccessTime time.Time `json:"access_time" example:"2027-03-03T08:15:00Z"`
}

// ProjectCreate godoc
// @Summary create project
// @Tags projects
// @Accept  json
// @Produce application/json
// @Param comment body CreateProjectRequest true "raw request body"
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
	var requestPayload CreateProjectRequest

	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422001, "could not read request body", err))
		return
	}

	if requestPayload.UserID == 0 {
		err := errors.New("userID does not be empty")
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422103, err.Error(), err))
		return
	}

	if requestPayload.Name == "" {
		err := errors.New("name does not be empty")
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422104, err.Error(), err))
		return
	}

	if requestPayload.AccessTime.IsZero() {
		err := errors.New("accessTime does not be empty")
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422104, err.Error(), err))
		return
	}

	if requestPayload.AccessTime.Before(time.Now()) {
		err := errors.New("accessTime should be more then NOW")
		_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422102, err.Error(), err))
		return
	}

	project, err := h.projectRepo.FindByUserID(requestPayload.UserID)
	if err != nil && err != gorm.ErrRecordNotFound {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500101, "could not compare new Project with projects", err))
		return
	}

	if project != nil {
		if project.Name == requestPayload.Name {
			err := fmt.Errorf("project with Name %s already exists", requestPayload.Name)
			_ = WriteResponse(w, http.StatusConflict, NewErrorPayload(409102, "project with Name already exists", err))
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
