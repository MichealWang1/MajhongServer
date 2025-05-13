package codes

import (
	"errors"
	"fmt"
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d]: %s", e.Code, e.Message)
}

func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

func NewWithCode(code int) *Error {
	return &Error{
		Code:    code,
		Message: GetMessage(code),
	}
}

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}

	return t.Code == e.Code
}

func ParseError(err error) *Error {
	e := &Error{}
	ok := errors.As(err, &e)
	if !ok {
		return &Error{
			Code:    UnKnowError,
			Message: GetMessage(UnKnowError),
		}
	}

	return e
}
