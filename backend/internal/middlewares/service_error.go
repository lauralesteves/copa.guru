package middlewares

import "fmt"

type ErrorCode int

const (
	ErrValidation   ErrorCode = 400
	ErrUnauthorized ErrorCode = 401
	ErrForbidden    ErrorCode = 403
	ErrNotFound     ErrorCode = 404
	ErrConflict     ErrorCode = 409
	ErrInternal     ErrorCode = 500
)

type ServiceError struct {
	Code    ErrorCode
	Message string
	Err     error
}

func (e *ServiceError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *ServiceError) Unwrap() error {
	return e.Err
}

func NewServiceError(code ErrorCode, message string, err error) *ServiceError {
	return &ServiceError{Code: code, Message: message, Err: err}
}

func NewValidationError(msg string) *ServiceError {
	return &ServiceError{Code: ErrValidation, Message: msg}
}

func NewUnauthorizedError(msg string) *ServiceError {
	return &ServiceError{Code: ErrUnauthorized, Message: msg}
}

func NewForbiddenError(msg string) *ServiceError {
	return &ServiceError{Code: ErrForbidden, Message: msg}
}

func NewNotFoundError(msg string) *ServiceError {
	return &ServiceError{Code: ErrNotFound, Message: msg}
}

func NewConflictError(msg string) *ServiceError {
	return &ServiceError{Code: ErrConflict, Message: msg}
}

func NewInternalError(msg string, err error) *ServiceError {
	return &ServiceError{Code: ErrInternal, Message: msg, Err: err}
}
