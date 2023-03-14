package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *BaseHandler) ProjectRead(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400001, "request is incorrect", err))
	}

	project, err := h.projectRepo.FindByID(projectID)
	if err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500001, "could not find Project by ProjectID", err))
	}

	payload := NewResponsePayload(fmt.Sprintf("project ID %d successfuly found", projectID), project)
	_ = WriteResponse(w, http.StatusOK, payload)
	return
}
