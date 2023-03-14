package handlers

import "net/http"

func (h *BaseHandler) ProjectDelete(w http.ResponseWriter, r *http.Request) {
	payload := NewResponsePayload("", "")
	_ = WriteResponse(w, http.StatusOK, payload)
	return
}
