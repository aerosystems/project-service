package CustomErrors

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

type EchoError struct {
	mode   EchoHandlerMode
	errors []ApiError
}

func NewEchoErrorHandler(mode string) echo.HTTPErrorHandler {
	e := EchoError{
		mode:   NewEchoHandlerMode(mode),
		errors: apiErrors,
	}
	return e.Handler
}

func (h *EchoError) Handler(err error, c echo.Context) {
	var code int
	var message map[string]interface{}
	var httpError *echo.HTTPError
	var apiError ApiError

	switch {
	case errors.As(err, &httpError):
		if httpError.Internal != nil {
			var internalError *echo.HTTPError
			if errors.As(httpError.Internal, &internalError) {
				code = internalError.Code
				message = map[string]interface{}{"message": internalError.Message}
			}
		} else {
			code = httpError.Code
			message = map[string]interface{}{"message": httpError.Message}
		}
	case errors.As(err, &apiError):
		code = apiError.HttpCode
		message = map[string]interface{}{"message": apiError.Message}
	default:
		code = http.StatusInternalServerError
		message = map[string]interface{}{"message": "Internal Server Error"}
		if h.mode == DevelopmentMode {
			message = map[string]interface{}{"message": err.Error()}
		}
	}

	if !c.Response().Committed {
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(httpError.Code)
		} else {
			err = c.JSON(code, message)
		}
		if err != nil {
			c.Echo().Logger.Error(err)
		}
	}
}
