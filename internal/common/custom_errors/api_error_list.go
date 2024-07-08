package CustomErrors

var apiErrors = []ApiError{
	ErrInvalidRequestBody,
	ErrInvalidRequestPayload,
	ErrForbidden,
	ErrProjectUuidInvalid,
	ErrProjectAlreadyExists,
	ErrProjectNameExists,
	ErrProjectNotFound,
	ErrProjectLimitExceeded,
}
