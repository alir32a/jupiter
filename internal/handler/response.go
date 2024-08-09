package handler

import (
	"errors"
	"github.com/alir32a/jupiter/internal/errorext"
	clog "github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"net/http"
)

const InternalServerError = "internal server error"

type HTTPResponse struct {
	OK     bool `json:"ok"`
	Result any  `json:"result"`
}

type Error struct {
	Error string `json:"error"`
}

func NewSuccessHTTPResponse(ctx echo.Context, status int, data any) error {
	return ctx.JSON(status, HTTPResponse{
		OK:     true,
		Result: data,
	})
}

func NewBindingError(ctx echo.Context, err error) error {
	return ctx.JSON(http.StatusBadRequest, HTTPResponse{
		OK:     false,
		Result: Error{Error: err.Error()},
	})
}

func NewFailedHTTPResponse(ctx echo.Context, logger *clog.Logger, err error) error {
	extErr := &errorext.Error{}

	if errors.As(err, &extErr) {
		status := extErr.Status()
		if status == 0 {
			status = http.StatusInternalServerError
		}

		return ctx.JSON(status, HTTPResponse{
			OK:     false,
			Result: Error{Error: err.Error()},
		})
	}

	logger.Error(err.Error())

	return ctx.JSON(http.StatusInternalServerError, HTTPResponse{
		OK:     false,
		Result: Error{Error: InternalServerError},
	})
}
