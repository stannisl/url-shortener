package domain

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrModelNotFound = errors.New("model not found")
	ErrNotUnique     = errors.New("not unique")
	ErrInvalidInput  = errors.New("invalid input")
)

type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err.Error())
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func NewNotFoundError(msg string) *Error {
	return &Error{
		Code:    http.StatusNotFound,
		Message: msg,
		Err:     ErrModelNotFound,
	}
}

func NewConflictError(msg string) *Error {
	return &Error{
		Code:    http.StatusConflict,
		Message: msg,
		Err:     ErrNotUnique,
	}
}

func NewBadRequestError(msg string, err error) *Error {
	return &Error{
		Code:    http.StatusBadRequest,
		Message: msg,
		Err:     err,
	}
}

func NewInternalError(msg string, err error) *Error {
	return &Error{
		Code:    http.StatusInternalServerError,
		Message: msg,
		Err:     err,
	}
}
