package gerrors

import (
	"encoding/json"
	"net/http"
)

type GlobalError struct {
	Code        int
	ServiceName string
	Message     string
	InnerErr    error
	StatusCode  int
}

func (ge *GlobalError) Error() string {
	b, _ := json.Marshal(ge)
	return string(b)
}

func New(code int, statusCode int, message string, err ...error) error {
	if message == "" {
		message = http.StatusText(statusCode)
	}
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	return &GlobalError{
		Code: code,
		// ServiceName: serviceName,
		Message:    message,
		InnerErr:   e,
		StatusCode: statusCode,
	}
}

// BadRequest generates a 400 error.
func BadRequest(code int, message string, err ...error) error {
	if message == "" {
		message = http.StatusText(http.StatusBadRequest)
	}
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	return &GlobalError{
		Code: code,
		// ServiceName: serviceName,
		StatusCode: http.StatusBadRequest,
		Message:    message,
		InnerErr:   e,
	}
}

// Unauthorized generates a 401 error.
func Unauthorized(code int, message string, err ...error) error {
	if message == "" {
		message = http.StatusText(http.StatusUnauthorized)
	}
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	return &GlobalError{
		Code: code,
		// ServiceName: serviceName,
		Message:    message,
		StatusCode: http.StatusUnauthorized,
		InnerErr:   e,
	}
}

// Forbidden generates a 403 error.
func Forbidden(code int, message string, err ...error) error {
	if message == "" {
		message = http.StatusText(http.StatusForbidden)
	}
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	return &GlobalError{
		Code: code,
		// ServiceName: serviceName,
		Message:    message,
		StatusCode: http.StatusForbidden,
		InnerErr:   e,
	}
}

// NotFound generates a 404 error.
func NotFound(code int, message string, err ...error) error {
	if message == "" {
		message = http.StatusText(http.StatusNotFound)
	}
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	return &GlobalError{
		Code: code,
		// ServiceName: serviceName,
		Message:    message,
		StatusCode: http.StatusNotFound,
		InnerErr:   e,
	}
}

// InternalServerError generates a 500 error.
func InternalServerError(code int, message string, err ...error) error {
	if message == "" {
		message = http.StatusText(http.StatusInternalServerError)
	}
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	return &GlobalError{
		Code: code,
		// ServiceName: serviceName,
		Message:    message,
		StatusCode: http.StatusInternalServerError,
		InnerErr:   e,
	}
}

// Conflict generates a 409 error.
func Conflict(code int, message string, err ...error) error {
	if message == "" {
		message = http.StatusText(http.StatusConflict)
	}
	var e error
	if len(err) > 0 {
		e = err[0]
	}
	return &GlobalError{
		Code: code,
		// ServiceName: serviceName,
		Message:    message,
		StatusCode: http.StatusConflict,
		InnerErr:   e,
	}
}

type ValidateError map[string][]string

func (ve ValidateError) Error() string {
	b, _ := json.Marshal(ve)
	return string(b)
}

// UnprocessableEntity generates a 422 error.
func UnprocessableEntity(code int, ve ValidateError) error {
	return &GlobalError{
		Code: code,
		// ServiceName: serviceName,
		Message:    "The given data failed to pass validation.",
		StatusCode: http.StatusUnprocessableEntity,
		InnerErr:   ve,
	}
}
