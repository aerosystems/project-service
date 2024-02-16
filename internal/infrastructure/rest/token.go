package rest

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type TokenHandler struct {
	*BaseHandler
	projectUsecase ProjectUsecase
}

func NewTokenHandler(baseHandler *BaseHandler, projectUsecase ProjectUsecase) *TokenHandler {
	return &TokenHandler{
		BaseHandler:    baseHandler,
		projectUsecase: projectUsecase,
	}
}

// ValidateToken godoc
// @Summary validate token
// @Tags token
// @Accept  json
// @Produce application/json
// @Security X-Api-Key
// @Success 204 {object} Response
// @Failure 401 {object} ErrorResponse
// @Router /v1/token/validate [get]
func (th TokenHandler) ValidateToken(c echo.Context) error {
	token := c.Request().Header.Get("X-Api-Key")
	if th.projectUsecase.IsProjectExistByToken(token) {
		return th.ErrorResponse(c, http.StatusUnauthorized, "could not get Project by Token", nil)
	}
	return c.JSON(http.StatusNoContent, nil)
}
