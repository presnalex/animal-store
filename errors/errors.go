package errors

import "github.com/pkg/errors"

type InternalError struct {
	Err error
}

func (e *InternalError) Error() string {
	return e.Err.Error()
}

func NewInternalError(err error) error {
	return errors.WithStack(&InternalError{Err: err})
}

type BadRequestError struct {
	err error
}

func (e *BadRequestError) Error() string {
	return "bad request: " + e.err.Error()
}

func NewBadRequestError(err error) error {
	return errors.WithStack(&BadRequestError{err: err})
}

type NotFoundError struct {
	Err string
}

func (e *NotFoundError) Error() string {
	return e.Err
}

func NewNotFoundError(err string) error {
	return errors.WithStack(&NotFoundError{Err: err})
}
