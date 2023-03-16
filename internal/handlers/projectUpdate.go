package handlers

import (
	"errors"
	"fmt"
	"github.com/aerosystems/project-service/internal/models"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

func (h *BaseHandler) ProjectUpdate(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400002, "request path param should be integer", err))
		return
	}

	var requestPayload models.ProjectRequest
	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400001, "request payload is incorrect", err))
		return
	}

	project, err := h.projectRepo.FindByID(projectID)
	if err != nil && err != gorm.ErrRecordNotFound {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500101, "could not compare new Project with projects in storage", err))
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
			err := errors.New("claim Access Time should be more then NOW")
			_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400102, err.Error(), err))
			return
		}
		project.AccessTime = requestPayload.AccessTime
	}

	if err = h.projectRepo.Update(project); err != nil {
		// TODO have to make user projects audit
		if err == gorm.ErrDuplicatedKey {
			_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500105, "user does not have the same project names", err))
			return
		}
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500104, "could not update project", err))
		return
	}

	payload := NewResponsePayload("project successfully updated", project)
	_ = WriteResponse(w, http.StatusOK, payload)
	return
}
