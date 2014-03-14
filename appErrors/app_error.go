// The appErrors package provides support for application specific error handling
package appErrors

import (
	"fmt"
)

const (
	VALIDATION_ERROR_MSG  = "validation_error"
	VALIDATION_ERROR_CODE = 409

	UNAUTHORIZED_ERROR_MSG  = "invalid_credentials"
	UNAUTHORIZED_ERROR_CODE = 401

	APP_ERROR_MSG  = "application_error"
	APP_ERROR_CODE = 500

	CACHE_CONTROL_HEADER = "Cache-control"

	NETWORK_READ_ERROR_CODE = 598
	NETWORK_READ_ERROR_MSG  = "network_read_error"
)

type (
	// AppError implements an application error which
	// reguires an error string and code
	AppError struct {
		ErrorMsg string
		Code     int
	}
)

// Error returns the error message that is associated with the AppError object
func (appError *AppError) Error() string {
	return appError.ErrorMsg
}

// ErrorCode returns the error code that is associated with the AppError object
func (appError *AppError) ErrorCode() int {
	if appError.Code == 0 {
		return APP_ERROR_CODE
	}

	return appError.Code
}

// NewError creates an AppError object
func NewError(error string, code int) error {
	return &AppError{ErrorMsg: error, Code: code}
}

// NewValidationError creates an AppError object for a validation error
func NewValidationError(error string) error {
	return &AppError{ErrorMsg: error, Code: VALIDATION_ERROR_CODE}
}

// Errorf creates a validation AppError based on the formatted string
func Errorf(format string, a ...interface{}) error {
	return NewError(fmt.Sprintf(format, a...), VALIDATION_ERROR_CODE)
}
