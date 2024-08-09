package handler

import (
	"context"
	"github.com/alir32a/jupiter/internal/model"
	clog "github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"net/http"
)

type UserService interface {
	GetAllUsers(ctx context.Context, req model.GetAllUsersRequest) (model.GetAllUsersResponse, error)
	BanUser(ctx context.Context, userID int) error
	UnbanUser(ctx context.Context, userID int) error
}

type UserHandler struct {
	svc    UserService
	logger *clog.Logger
}

func NewUserHandler(svc UserService, logger *clog.Logger) *UserHandler {
	return &UserHandler{svc: svc, logger: logger}
}

func (u UserHandler) GetAllUsers(ctx echo.Context) error {
	var req GetAllUsersRequest

	if err := ctx.Bind(&req); err != nil {
		return NewBindingError(ctx, err)
	}

	resp, err := u.svc.GetAllUsers(ctx.Request().Context(), toModelGetAllUsersRequest(req))
	if err != nil {
		return NewFailedHTTPResponse(ctx, u.logger, err)
	}

	return NewSuccessHTTPResponse(ctx, http.StatusOK, toCtrlGetAllUsersResponse(resp))
}

func (u UserHandler) BanUser(ctx echo.Context) error {
	var req BanUserRequest

	if err := ctx.Bind(&req); err != nil {
		return NewBindingError(ctx, err)
	}

	err := u.svc.BanUser(ctx.Request().Context(), req.ID)
	if err != nil {
		return NewFailedHTTPResponse(ctx, u.logger, err)
	}

	return NewSuccessHTTPResponse(ctx, http.StatusOK, nil)
}

func (u UserHandler) UnBanUser(ctx echo.Context) error {
	var req BanUserRequest

	if err := ctx.Bind(&req); err != nil {
		return NewBindingError(ctx, err)
	}

	err := u.svc.UnbanUser(ctx.Request().Context(), req.ID)
	if err != nil {
		return NewFailedHTTPResponse(ctx, u.logger, err)
	}

	return NewSuccessHTTPResponse(ctx, http.StatusOK, nil)
}

func (u UserHandler) SetRoutes(router *echo.Group) {
	router.GET("/users", u.GetAllUsers)
	router.POST("/users/:id/ban", u.BanUser)
	router.POST("/users/:id/unban", u.UnBanUser)
}
