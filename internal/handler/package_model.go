package handler

import (
	"github.com/alir32a/jupiter/internal/model"
	"github.com/alir32a/jupiter/pkg/util"
	"time"
)

type PackageEntity struct {
	ID                   int        `json:"id"`
	UserID               int        `json:"user_id"`
	Username             string     `json:"username"`
	TrafficLimit         string     `json:"traffic_limit"`
	DownloadTrafficUsage string     `json:"download_traffic_usage"`
	UploadTrafficUsage   string     `json:"upload_traffic_usage"`
	MaxConnections       int        `json:"max_connections"`
	IsTrial              bool       `json:"is_trial"`
	ExpirationInDays     int        `json:"expiration_in_days"`
	ExpireAt             *time.Time `json:"expire_at"`
	CreatedAt            time.Time  `json:"created_at"`
}

type GetPackagesRequest struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Username string `query:"username"`
}

type GetPackagesResponse struct {
	Pagination
	Packages []PackageEntity `json:"packages"`
}

type CreatePackageRequest struct {
	Username       string `json:"username"`
	TrafficLimit   int    `json:"traffic_limit"`
	MaxConnections int    `json:"max_connections"`
	Expiry         int    `json:"expiry"`
}

type GetUserActiveAndReservedPackagesRequest struct {
	UserID int `query:"user_id"`
}

type PackageSummary struct {
	Status       string     `json:"status"`
	TrafficLimit string     `json:"traffic_limit"`
	UsedTraffic  string     `json:"used_traffic"`
	ExpiresAt    *time.Time `json:"expires_at"`
}

func toModelGetPackagesRequest(req GetPackagesRequest) model.GetPackagesRequest {
	return model.GetPackagesRequest{
		Pagination: model.Pagination{
			CurrentPage: req.Page,
			PageSize:    req.PageSize,
		},
		Username: req.Username,
	}
}

func toCtrlPackageEntity(req model.PackageEntity) PackageEntity {
	return PackageEntity{
		ID:                   req.ID,
		UserID:               req.UserID,
		Username:             req.Username,
		TrafficLimit:         util.ToHumanReadableBytes(req.TrafficLimit),
		DownloadTrafficUsage: util.ToHumanReadableBytes(req.DownloadTrafficUsage),
		UploadTrafficUsage:   util.ToHumanReadableBytes(req.UploadTrafficUsage),
		MaxConnections:       req.MaxConnections,
		IsTrial:              req.IsTrial,
		ExpirationInDays:     req.ExpirationInDays,
		ExpireAt:             req.ExpireAt,
		CreatedAt:            req.CreatedAt,
	}
}

func toCtrlPackageEntities(packages []model.PackageEntity) []PackageEntity {
	result := make([]PackageEntity, 0, len(packages))

	for _, pack := range packages {
		result = append(result, toCtrlPackageEntity(pack))
	}

	return result
}

func toModelCreatePackageRequest(req CreatePackageRequest) model.CreatePackageRequest {
	return model.CreatePackageRequest{
		Username:         req.Username,
		Traffic:          req.TrafficLimit,
		MaxConnections:   req.MaxConnections,
		ExpirationInDays: req.Expiry,
	}
}

func toPackageSummary(status string, pack model.PackageEntity) PackageSummary {
	return PackageSummary{
		Status:       status,
		TrafficLimit: util.ToHumanReadableBytes(pack.TrafficLimit),
		UsedTraffic:  util.ToHumanReadableBytes(pack.DownloadTrafficUsage + pack.UploadTrafficUsage),
		ExpiresAt:    pack.ExpireAt,
	}
}

func toCtrlGetUserActiveAndReservedPackagesResponse(req model.GetUserPackages) []PackageSummary {
	result := []PackageSummary{toPackageSummary("active", req.ActivePackage)}

	for _, pack := range req.ReservedPackages {
		result = append(result, toPackageSummary("reserved", pack))
	}

	return result
}
