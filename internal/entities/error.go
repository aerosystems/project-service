package entities

import (
	"github.com/aerosystems/common-service/customerrors"
	"google.golang.org/grpc/codes"
	"net/http"
)

var (
	ErrInvalidRequestBody    = customerrors.InternalError{Message: "Invalid request body.", HttpCode: http.StatusUnprocessableEntity, GrpcCode: codes.InvalidArgument}
	ErrInvalidRequestPayload = customerrors.InternalError{Message: "Invalid request payload.", HttpCode: http.StatusBadRequest, GrpcCode: codes.InvalidArgument}
	ErrForbidden             = customerrors.InternalError{Message: "Forbidden.", HttpCode: http.StatusForbidden, GrpcCode: codes.PermissionDenied}
	ErrUnknownUserRole       = customerrors.InternalError{Message: "Unknown user role.", HttpCode: http.StatusForbidden, GrpcCode: codes.PermissionDenied}
	ErrProjectUuidInvalid    = customerrors.InternalError{Message: "Project UUID is invalid.", HttpCode: http.StatusBadRequest, GrpcCode: codes.InvalidArgument}
	ErrProjectAlreadyExists  = customerrors.InternalError{Message: "Project already exists.", HttpCode: http.StatusConflict, GrpcCode: codes.AlreadyExists}
	ErrProjectNameExists     = customerrors.InternalError{Message: "Project name already exists.", HttpCode: http.StatusConflict, GrpcCode: codes.AlreadyExists}
	ErrProjectNotFound       = customerrors.InternalError{Message: "Project not found.", HttpCode: http.StatusNotFound, GrpcCode: codes.NotFound}
	ErrProjectLimitExceeded  = customerrors.InternalError{Message: "Project limit exceeded.", HttpCode: http.StatusForbidden, GrpcCode: codes.PermissionDenied}
)
