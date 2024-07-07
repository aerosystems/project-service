package CustomErrors

import (
	"net/http"
)

type ApiError struct {
	Message  string
	HttpCode int
}

func (e ApiError) Error() string {
	return e.Message
}

var (
	ErrInvalidRequestBody    = ApiError{"Invalid request body.", http.StatusUnprocessableEntity}
	ErrInvalidRequestPayload = ApiError{"Invalid request payload.", http.StatusBadRequest}
	ErrForbidden             = ApiError{"Forbidden.", http.StatusForbidden}
	ErrProjectUuidInvalid    = ApiError{"Project UUID is invalid.", http.StatusBadRequest}
	ErrProjectAlreadyExists  = ApiError{"Project already exists.", http.StatusConflict}
	ErrProjectNameExists     = ApiError{"Project name already exists.", http.StatusConflict}
	ErrProjectNotFound       = ApiError{"Project not found.", http.StatusNotFound}
	ErrProjectLimitExceeded  = ApiError{"Project limit exceeded.", http.StatusForbidden}
)
