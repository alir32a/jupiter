package handler

import (
	"context"
	"errors"
	"github.com/alir32a/jupiter/config"
	"github.com/alir32a/jupiter/internal/model"
	"github.com/alir32a/jupiter/pkg/jwt"
	clog "github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type AdminService interface {
	Login(ctx context.Context, req model.AdminLoginRequest) error
	ChangePassword(ctx context.Context, req model.ChangePasswordRequest) error
}

type AdminHandler struct {
	svc    AdminService
	cfg    *config.HTTPServerConfig
	logger *clog.Logger
}

func NewAdminHandler(svc AdminService, cfg *config.HTTPServerConfig, logger *clog.Logger) *AdminHandler {
	return &AdminHandler{
		svc:    svc,
		cfg:    cfg,
		logger: logger,
	}
}

func (a AdminHandler) Login(ctx echo.Context) error {
	var req AdminLoginRequest

	if err := ctx.Bind(&req); err != nil {
		return NewBindingError(ctx, err)
	}

	err := a.svc.Login(ctx.Request().Context(), toModelAdminLoginRequest(req))
	if err != nil {
		return NewFailedHTTPResponse(ctx, a.logger, err)
	}

	token, err := jwt.CreateToken(req.Username, a.cfg.AccessTokenSecret, a.cfg.AccessTokenExpireTime)
	if err != nil {
		return NewFailedHTTPResponse(ctx, a.logger, err)
	}

	cookie := &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    token,
		Expires:  time.Now().Add(a.cfg.AccessTokenExpireTime),
		MaxAge:   int(a.cfg.AccessTokenExpireTime.Seconds()),
		HttpOnly: true,
		Secure:   true,
	}

	if a.cfg.ENV == "debug" {
		cookie.SameSite = http.SameSiteNoneMode
	}

	ctx.SetCookie(cookie)

	return NewSuccessHTTPResponse(ctx, http.StatusOK, nil)
}

func (a AdminHandler) Self(ctx echo.Context) error {
	cookie, err := ctx.Cookie(AccessTokenCookieName)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, nil)
	}

	claim, err := jwt.ParseToken(cookie.Value, a.cfg.AccessTokenSecret)
	if err != nil {
		return NewFailedHTTPResponse(ctx, a.logger, err)
	}

	return NewSuccessHTTPResponse(ctx, http.StatusOK, SelfResponse{Username: claim.Username})
}

func (a AdminHandler) Logout(ctx echo.Context) error {
	ctx.SetCookie(&http.Cookie{Name: AccessTokenCookieName, MaxAge: -1})

	return NewSuccessHTTPResponse(ctx, http.StatusOK, nil)
}

func (a AdminHandler) ChangePassword(ctx echo.Context) error {
	var req ChangePasswordRequest

	if err := ctx.Bind(&req); err != nil {
		return NewBindingError(ctx, err)
	}

	if req.NewPassword != req.ConfirmPassword {
		return NewBindingError(ctx, errors.New("passwords are not equal"))
	}

	claim, ok := ctx.Get("user").(*jwt.Claim)
	if !ok {
		return NewFailedHTTPResponse(ctx, a.logger, errors.New("user not found"))
	}
	req.Username = claim.Username

	err := a.svc.ChangePassword(ctx.Request().Context(), toModelChangePasswordRequest(req))
	if err != nil {
		return NewFailedHTTPResponse(ctx, a.logger, err)
	}

	return a.Logout(ctx)
}

func (a AdminHandler) SetNoAuthRoutes(router *echo.Group) {
	router.POST("/login", a.Login)
}

func (a AdminHandler) SetRoutes(router *echo.Group) {
	router.GET("/self", a.Self)
	router.POST("/logout", a.Logout)
	router.POST("/change-password", a.ChangePassword)
}
