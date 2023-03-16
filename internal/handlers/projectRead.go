package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (h *BaseHandler) ProjectRead(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400002, "request path param should be integer", err))
		return
	}

	project, err := h.projectRepo.FindByID(projectID)
	if err != nil && err != gorm.ErrRecordNotFound {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500001, "could not find Project by Project ID", err))
		return
	}

	if project == nil {
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404001, "project not found", err))
		return
	}

	payload := NewResponsePayload(fmt.Sprintf("project ID %d successfuly found", projectID), project)
	_ = WriteResponse(w, http.StatusOK, payload)
	return
}
