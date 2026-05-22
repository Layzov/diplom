package apperror

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrValidation  = errors.New("validation error")
	ErrConflict    = errors.New("conflict")
)

// ValidationError carries a client-safe message.
type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) Is(target error) bool {
	return target == ErrValidation
}

func Validation(msg string) error {
	return &ValidationError{Message: msg}
}
