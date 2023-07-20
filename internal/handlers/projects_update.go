package handlers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type UpdateProjectRequest struct {
	Name       string    `json:"name" example:"bla-bla-bla.com"`
	AccessTime time.Time `json:"access_time" example:"2027-03-03T08:15:00Z"`
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

	project, err := h.projectRepo.FindByID(projectID)
	if err != nil && err != gorm.ErrRecordNotFound {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500101, "could not compare new Project with projects", err))
		return
	}
	if project == nil {
		err := fmt.Errorf("project ID %d does not exist", projectID)
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404005, err.Error(), err))
		return
	}

	if requestPayload.Name != "" {
		project.Name = requestPayload.Name
	}

	if !requestPayload.AccessTime.IsZero() {
		if requestPayload.AccessTime.Before(time.Now()) {
			err := errors.New("accessTime should be more then NOW")
			_ = WriteResponse(w, http.StatusUnprocessableEntity, NewErrorPayload(422102, err.Error(), err))
			return
		}
		project.AccessTime = requestPayload.AccessTime
	}

	if err = h.projectRepo.Update(project); err != nil {
		// TODO have to make user projects audit
		if err == gorm.ErrDuplicatedKey {
			_ = WriteResponse(w, http.StatusConflict, NewErrorPayload(500105, "user does not have the same project names", err))
			return
		}
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500104, "could not update project", err))
		return
	}

	payload := NewResponsePayload("project successfully updated", project)
	_ = WriteResponse(w, http.StatusOK, payload)
	return
}
