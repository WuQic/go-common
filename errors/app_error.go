package errors

const (
	VALIDATION_ERROR_MSG  = "Validation Error: The parameters provided were invalid."
	VALIDATION_ERROR_CODE = 409

	UNAUTHORIZED_ERROR_MSG  = "Invalid Credentials were supplied."
	UNAUTHORIZED_ERROR_CODE = 401

	APP_ERROR_MSG  = "An Application Error has occured."
	APP_ERROR_CODE = 500

	CACHE_CONTROL_HEADER = "Cache-control"
)

type (
	AppError struct {
		ErrorMsg string
		Code     int
	}
)

func (this *AppError) Error() string {
	return this.ErrorMsg
}

func (this *AppError) ErrorCode() int {
	if this.Code == 0 {
		return APP_ERROR_CODE
	}

	return this.Code
}

func NewError(error string, code int) error {
	return &AppError{ErrorMsg: error, Code: code}
}

func NewValidationError(error string) error {
	return &AppError{ErrorMsg: error, Code: VALIDATION_ERROR_CODE}
}
