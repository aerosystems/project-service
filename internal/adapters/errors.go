package adapters

import (
	CustomErrors "github.com/aerosystems/project-service/internal/common/custom_errors"
	"net/http"
)

var (
	ErrProjectNotFound = CustomErrors.ApiError{Message: "Project not found.", HttpCode: http.StatusNotFound}
)
