package handlers

import (
	"github.com/labstack/echo/v4"
)

// ValidateToken godoc
// @Summary validate token
// @Tags token
// @Accept  json
// @Produce application/json
// @Security X-Api-Key
// @Success 204 {object} Response
// @Failure 401 {object} ErrorResponse
// @Router /v1/token/validate [get]
func (h *BaseHandler) ValidateToken(c echo.Context) error {
	//project, ok := r.Context().Value(helpers.ContextKey("project")).(*models.Project)
	//if !ok {
	//	err := fmt.Errorf("could not get Project by Token: %s", r.Header.Get("X-Api-Key"))
	//	_ = WriteResponse(w, http.StatusUnauthorized, NewErrorPayload(401001, "could not get Project by Token", err))
	//	return
	//}
	//
	//if project == nil {
	//	err := fmt.Errorf("could not get Project by Token: %s", r.Header.Get("X-Api-Key"))
	//	_ = WriteResponse(w, http.StatusUnauthorized, NewErrorPayload(401002, "could not get Project by Token", err))
	//	return
	//}
	//
	//_ = WriteResponse(w, http.StatusNoContent, nil)
	return nil
}
