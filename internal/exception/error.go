package exception

import (
	"errors"
	"fmt"
)

type AppError struct {
	message string
}

func (e *AppError) Error() string {
	return e.message
}

func New(text string) *AppError {
	return &AppError{message: text}
}

func Wrap(err, wrapped error) error {
	return fmt.Errorf("%w: %w", err, wrapped)
}

func Is(err error) bool {
	var appError *AppError
	return err != nil && errors.As(err, &appError)
}
