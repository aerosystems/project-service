package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func (h *BaseHandler) ProjectDelete(w http.ResponseWriter, r *http.Request) {
	projectID, err := strconv.Atoi(chi.URLParam(r, "projectID"))
	if err != nil {
		_ = WriteResponse(w, http.StatusBadRequest, NewErrorPayload(400002, "request path param should be integer", err))
		return
	}

	project, err := h.projectRepo.FindByID(projectID)

	if err != nil && err != gorm.ErrRecordNotFound {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500001, "could not find Project by ProjectID", err))
		return
	}

	if project == nil {
		_ = WriteResponse(w, http.StatusNotFound, NewErrorPayload(404001, "project not found", err))
		return
	}

	if err = h.projectRepo.Delete(project); err != nil {
		_ = WriteResponse(w, http.StatusInternalServerError, NewErrorPayload(500106, "could not delete Project", err))
		return
	}

	payload := NewResponsePayload(fmt.Sprintf("project ID %d successfuly deleted", projectID), nil)
	_ = WriteResponse(w, http.StatusOK, payload)
	return
}
