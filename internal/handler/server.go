package handler

import (
	"fmt"
	"github.com/alir32a/jupiter/config"
	clog "github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
)

type HTTPServer struct {
	*echo.Echo
	cfg    *config.HTTPServerConfig
	logger *clog.Logger
}

func NewHTTPServer(cfg *config.HTTPServerConfig, logger *clog.Logger) *HTTPServer {
	return &HTTPServer{
		Echo:   echo.New(),
		cfg:    cfg,
		logger: logger,
	}
}

func (h *HTTPServer) Run() error {
	return h.Start(fmt.Sprintf("%s:%d", h.cfg.Host, h.cfg.Port))
}
