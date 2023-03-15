package handlers

import (
	"fmt"
	"github.com/aerosystems/project-service/internal/models"
	"gorm.io/gorm"
	"net/http"
)

func (h *BaseHandler) ProjectCreate(w http.ResponseWriter, r *http.Request) {
	var requestPayload models.ProjectRequest

	if err := ReadRequest(w, r, &requestPayload); err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400001, "request payload is incorrect", err))
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
