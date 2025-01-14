package HTTPServer

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type TokenHandler struct {
	*BaseHandler
	tokenUsecase TokenUsecase
}

func NewTokenHandler(baseHandler *BaseHandler, tokenUsecase TokenUsecase) *TokenHandler {
	return &TokenHandler{
		BaseHandler:  baseHandler,
		tokenUsecase: tokenUsecase,
	}
}

// ValidateToken godoc
// @Summary validate token
// @Tags token
// @Accept  json
// @Produce application/json
// @Security X-Api-Key
// @Success 204
// @Failure 401 {object} ErrorResponse
// @Router /v1/token/validate [get]
func (th TokenHandler) ValidateToken(c echo.Context) error {
	token := c.Request().Header.Get("X-Api-Key")
	if !th.tokenUsecase.IsProjectExistByToken(token) {
		return echo.NewHTTPError(http.StatusUnauthorized, "could not get Project by Token")
	}

	return c.JSON(http.StatusNoContent, nil)
}
