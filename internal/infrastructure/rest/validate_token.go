package rest

import (
	"github.com/labstack/echo/v4"
	"net/http"
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
	token := c.Request().Header.Get("X-Api-Key")
	if h.projectService.IsProjectExistByToken(token) {
		return h.ErrorResponse(c, http.StatusUnauthorized, "could not get Project by Token", nil)
	}
	return c.JSON(http.StatusNoContent, nil)
}
