package apperror

import "fmt"

type AppError struct {
	Code       int
	StatusCode int
	Message    string
	Errors     []string
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

func New(statusCode, code int, message string, errors []string) *AppError {
	return &AppError{
		Code:       code,
		StatusCode: statusCode,
		Message:    message,
		Errors:     errors,
	}
}

// Predefined Errors
var (
	ErrBadRequest      = New(400, 400000, "bad request", nil)
	ErrEmailRegistered = New(400, 400000, "email already registered", nil)
	ErrUserNotFound    = New(404, 404000, "user not found", nil)
	ErrWrongPassword   = New(401, 401008, "wrong password", nil)
	ErrUnauthorized    = New(401, 401000, "unauthorized", nil)
	ErrInvalidToken    = New(401, 401001, "invalid or expired token", nil)
	ErrInternalServer  = New(500, 500000, "internal server error", nil)
)
