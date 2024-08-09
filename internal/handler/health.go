package handler

import (
	clog "github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"net/http"
)

type HealthCheckHandler struct {
	db     *gorm.DB
	logger *clog.Logger
}

func NewHealthCheckHandler(db *gorm.DB) *HealthCheckHandler {
	return &HealthCheckHandler{db: db}
}

func (h HealthCheckHandler) HealthCheck(ctx echo.Context) error {
	sqlDb, err := h.db.DB()
	if err != nil {
		return NewFailedHTTPResponse(ctx, h.logger, err)
	}

	if err := sqlDb.Ping(); err != nil {
		return NewFailedHTTPResponse(ctx, h.logger, err)
	}

	return NewSuccessHTTPResponse(ctx, http.StatusOK, nil)
}

func (h HealthCheckHandler) SetRoutes(router *echo.Group) {
	router.GET("/health", h.HealthCheck)
}
