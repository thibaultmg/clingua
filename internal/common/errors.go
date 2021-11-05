package common

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidData   = errors.New("invalid data")
	//ErrInternalError = errors.New("internal error")
)

type ErrInternalError struct {
	wrapped error
}

func (e ErrInternalError) Error() string {
	return fmt.Sprintf("internal error: %v", e.wrapped)
}

func (e ErrInternalError) Unwrap() error {
	return e.wrapped
}

func NewErrInternalError(err error) ErrInternalError {
	return ErrInternalError{
		wrapped: err,
	}
}
