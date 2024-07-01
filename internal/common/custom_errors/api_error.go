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
	ErrProjectUuidInvalid     = ApiError{"Project UUID is invalid.", http.StatusBadRequest}
	ErrProjectAlreadyExists   = ApiError{"Project already exists.", http.StatusConflict}
	ErrProjectNameExists      = ApiError{"Project name already exists.", http.StatusConflict}
	ErrProjectNotFound        = ApiError{"Project not found.", http.StatusNotFound}
	ErrProjectDeleteForbidden = ApiError{"Project delete forbidden.", http.StatusForbidden}
	ErrProjectUpdateForbidden = ApiError{"Project update forbidden.", http.StatusForbidden}
	ErrProjectLimitExceeded   = ApiError{"Project limit exceeded.", http.StatusForbidden}
)
