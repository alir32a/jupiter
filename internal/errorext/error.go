package errorext

import (
	clog "github.com/charmbracelet/log"
	"net/http"
	"runtime/debug"
)

type Error struct {
	message    string
	status     int
	stackTrace []byte
}

func (e Error) Error() string {
	return e.message
}

func (e Error) Status() int {
	return e.status
}

func New(msg string) *Error {
	return &Error{
		message:    msg,
		stackTrace: debug.Stack(),
	}
}

func NewInternalError(logger *clog.Logger, err error) error {
	logger.Error(err.Error())

	return &Error{
		message:    err.Error(),
		status:     http.StatusInternalServerError,
		stackTrace: debug.Stack(),
	}
}

func NewNotFoundError(err error) error {
	return &Error{
		message:    err.Error(),
		status:     http.StatusNotFound,
		stackTrace: debug.Stack(),
	}
}

func NewBadRequestError(err error) error {
	return &Error{
		message:    err.Error(),
		status:     http.StatusBadRequest,
		stackTrace: debug.Stack(),
	}
}
