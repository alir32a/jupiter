package handler

import (
	"context"
	"github.com/alir32a/jupiter/internal/model"
	clog "github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"net/http"
)

type ConnectionService interface {
	GetActiveConnections(ctx context.Context, req model.GetActiveConnectionsRequest) (model.GetActiveConnectionsResponse, error)
	GetSystemStatus(ctx context.Context) (model.GetSystemStatusResponse, error)
	DisconnectID(ctx context.Context, id int) error
	GetUserActiveConnections(ctx context.Context, username string) ([]model.ConnectionEntity, error)
}

type ConnectionHandler struct {
	svc    ConnectionService
	logger *clog.Logger
}

func NewConnectionHandler(svc ConnectionService, logger *clog.Logger) *ConnectionHandler {
	return &ConnectionHandler{svc: svc, logger: logger}
}

func (c ConnectionHandler) GetActiveConnections(ctx echo.Context) error {
	var req GetActiveConnectionsRequest

	if err := ctx.Bind(&req); err != nil {
		return NewBindingError(ctx, err)
	}

	resp, err := c.svc.GetActiveConnections(ctx.Request().Context(), toModelGetActiveConnectionsRequest(req))
	if err != nil {
		return NewFailedHTTPResponse(ctx, c.logger, err)
	}

	return NewSuccessHTTPResponse(ctx, http.StatusOK, GetActiveConnectionsResponse{
		Connections: toCtrlConnectionEntities(resp.Connections),
		Pagination:  toCtrlPagination(*resp.Pagination),
	})
}

func (c ConnectionHandler) GetSystemStatus(ctx echo.Context) error {
	resp, err := c.svc.GetSystemStatus(ctx.Request().Context())
	if err != nil {
		return NewFailedHTTPResponse(ctx, c.logger, err)
	}

	return NewSuccessHTTPResponse(ctx, http.StatusOK, toCtrlGetSystemStatusResponse(resp))
}

func (c ConnectionHandler) DisconnectID(ctx echo.Context) error {
	var req DisconnectIDRequest

	if err := ctx.Bind(&req); err != nil {
		return NewBindingError(ctx, err)
	}

	err := c.svc.DisconnectID(ctx.Request().Context(), req.ID)
	if err != nil {
		return NewFailedHTTPResponse(ctx, c.logger, err)
	}

	return NewSuccessHTTPResponse(ctx, http.StatusOK, nil)
}

func (c ConnectionHandler) GetUserActiveConnections(ctx echo.Context) error {
	var req GetUserActiveConnectionsRequest

	if err := ctx.Bind(&req); err != nil {
		return NewBindingError(ctx, err)
	}

	connections, err := c.svc.GetUserActiveConnections(ctx.Request().Context(), req.Username)
	if err != nil {
		return NewFailedHTTPResponse(ctx, c.logger, err)
	}

	return NewSuccessHTTPResponse(ctx, http.StatusOK, toCtrlConnectionSummaries(connections))
}

func (c ConnectionHandler) SetRoutes(router *echo.Group) {
	router.GET("/connections", c.GetActiveConnections)
	router.GET("/system-statuses", c.GetSystemStatus)
	router.POST("/connections/:id/disconnect", c.DisconnectID)
	router.GET("/user-connections", c.GetUserActiveConnections)
}
