package token

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
// @Success 204 {object}
// @Failure 401 {object} echo.HTTPError
// @Router /v1/token/validate [get]
func (th Handler) ValidateToken(c echo.Context) error {
	token := c.Request().Header.Get("X-Api-Key")
	if !th.tokenUsecase.IsProjectExistByToken(token) {
		return echo.NewHTTPError(http.StatusUnauthorized, "could not get Project by Token")
	}
	return c.JSON(http.StatusNoContent, nil)
}
