package handler

import (
	"context"
	"github.com/alir32a/jupiter/internal/model"
	clog "github.com/charmbracelet/log"
	"github.com/labstack/echo/v4"
	"net/http"
)

type PackageService interface {
	GetPackages(ctx context.Context, req model.GetPackagesRequest) (model.GetPackagesResponse, error)
	CreatePackage(ctx context.Context, req model.CreatePackageRequest) error
	GetUserActiveAndReservedPackages(ctx context.Context, userID int) (model.GetUserPackages, error)
}

type PackageHandler struct {
	svc    PackageService
	logger *clog.Logger
}

func NewPackageHandler(svc PackageService, logger *clog.Logger) *PackageHandler {
	return &PackageHandler{svc: svc, logger: logger}
}

func (p PackageHandler) GetPackages(ctx echo.Context) error {
	var req GetPackagesRequest

	if err := ctx.Bind(&req); err != nil {
		return NewBindingError(ctx, err)
	}

	packages, err := p.svc.GetPackages(ctx.Request().Context(), toModelGetPackagesRequest(req))
	if err != nil {
		return NewFailedHTTPResponse(ctx, p.logger, err)
	}

	return NewSuccessHTTPResponse(ctx, http.StatusOK, GetPackagesResponse{
		Pagination: toCtrlPagination(packages.Pagination),
		Packages:   toCtrlPackageEntities(packages.Packages),
	})
}

func (p PackageHandler) CreatePackage(ctx echo.Context) error {
	var req CreatePackageRequest

	if err := ctx.Bind(&req); err != nil {
		return NewBindingError(ctx, err)
	}

	err := p.svc.CreatePackage(ctx.Request().Context(), toModelCreatePackageRequest(req))
	if err != nil {
		return NewFailedHTTPResponse(ctx, p.logger, err)
	}

	return NewSuccessHTTPResponse(ctx, http.StatusCreated, nil)
}

func (p PackageHandler) GetUserActiveAndReservedPackages(ctx echo.Context) error {
	var req GetUserActiveAndReservedPackagesRequest

	if err := ctx.Bind(&req); err != nil {
		return NewBindingError(ctx, err)
	}

	resp, err := p.svc.GetUserActiveAndReservedPackages(ctx.Request().Context(), req.UserID)
	if err != nil {
		return NewFailedHTTPResponse(ctx, p.logger, err)
	}

	return NewSuccessHTTPResponse(ctx, http.StatusOK, toCtrlGetUserActiveAndReservedPackagesResponse(resp))
}

func (p PackageHandler) SetRoutes(router *echo.Group) {
	router.GET("/packages", p.GetPackages)
	router.POST("/packages", p.CreatePackage)
	router.GET("/active-packages", p.GetUserActiveAndReservedPackages)
}
