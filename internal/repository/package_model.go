package repository

import (
	"github.com/alir32a/jupiter/internal/model"
	"time"
)

type PackageEntity struct {
	ID                   int
	UserID               int
	TrafficLimit         int
	DownloadTrafficUsage int
	UploadTrafficUsage   int
	MaxConnections       int
	IsTrial              bool
	ExpirationInDays     int
	ExpireAt             *time.Time
	CreatedAt            time.Time
}

func (PackageEntity) TableName() string {
	return "package"
}

func toModelPackageEntity(req PackageEntity) model.PackageEntity {
	return model.PackageEntity{
		ID:                   req.ID,
		UserID:               req.UserID,
		TrafficLimit:         req.TrafficLimit,
		DownloadTrafficUsage: req.DownloadTrafficUsage,
		UploadTrafficUsage:   req.UploadTrafficUsage,
		MaxConnections:       req.MaxConnections,
		ExpirationInDays:     req.ExpirationInDays,
		ExpireAt:             req.ExpireAt,
		CreatedAt:            req.CreatedAt,
	}
}

func toModelPackageEntities(req []PackageEntity) []model.PackageEntity {
	result := make([]model.PackageEntity, 0, len(req))

	for _, pack := range req {
		result = append(result, toModelPackageEntity(pack))
	}

	return result
}
