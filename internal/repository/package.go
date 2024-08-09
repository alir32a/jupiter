package repository

import (
	"context"
	"errors"
	"github.com/alir32a/jupiter/internal/model"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
	"time"
)

type PackageRepository struct {
	db *gorm.DB
}

func NewPackageRepository(db *gorm.DB) *PackageRepository {
	return &PackageRepository{db: db}
}

func (p PackageRepository) CreatePackage(ctx context.Context, req model.CreatePackageRequest) error {
	pack := PackageEntity{
		UserID:           req.UserID,
		TrafficLimit:     req.Traffic,
		MaxConnections:   req.MaxConnections,
		IsTrial:          req.IsTrial,
		ExpirationInDays: req.ExpirationInDays,
	}

	activePack, err := p.getUserActivePackage(ctx, req.UserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if activePack.ID != 0 {
		expireAt := time.Now().Add(time.Duration(req.ExpirationInDays) * 24 * time.Hour)

		pack.ExpireAt = &expireAt
	}

	return p.db.WithContext(ctx).Create(&pack).Error
}

func (p PackageRepository) GetUsersActiveAndReservedPackages(ctx context.Context, userIDs ...int) (model.GetUsersActivePackagesResponse, error) {
	var (
		packages []PackageEntity
		packMap  = map[int]model.GetUserPackages{}
	)

	err := p.db.
		WithContext(ctx).
		Raw(`with packages as (select * from package 
			 where download_traffic_usage + upload_traffic_usage < traffic_limit and (expire_at > now() or expire_at is null)) 
			 select * from packages where user_id in ?`, userIDs).Scan(&packages).Error
	if err != nil {
		return model.GetUsersActivePackagesResponse{}, err
	}

	for _, pack := range packages {
		userPacks := packMap[pack.UserID]
		userPacks.UserID = pack.UserID

		if pack.ExpireAt != nil {
			userPacks.ActivePackage = toModelPackageEntity(pack)

			continue
		}

		userPacks.ReservedPackages = append(userPacks.ReservedPackages, toModelPackageEntity(pack))
		packMap[pack.UserID] = userPacks
	}

	return model.GetUsersActivePackagesResponse{Packages: maps.Values(packMap)}, nil
}

func (p PackageRepository) GetUserActiveAndReservedPackages(ctx context.Context, userID int) (model.GetUserPackages, error) {
	var (
		packages []PackageEntity
		result   model.GetUserPackages
	)

	err := p.db.
		WithContext(ctx).
		Model(&PackageEntity{}).
		Where("download_traffic_usage + upload_traffic_usage < traffic_limit and (expire_at > now() or expire_at is null)").
		Where("user_id = ?", userID).
		Scan(&packages).Error
	if err != nil {
		return model.GetUserPackages{}, err
	}

	for _, pack := range packages {
		result.UserID = pack.UserID

		if pack.ExpireAt != nil {
			result.ActivePackage = toModelPackageEntity(pack)

			continue
		}

		result.ReservedPackages = append(result.ReservedPackages, toModelPackageEntity(pack))
	}

	return result, nil
}

func (p PackageRepository) getUserActivePackage(ctx context.Context, userID int) (PackageEntity, error) {
	var pack PackageEntity

	return pack, p.db.
		WithContext(ctx).
		Model(&PackageEntity{}).
		Where("user_id = ?", userID).
		Where("download_traffic_usage + upload_traffic_usage < traffic_limit and expire_at > now()").
		First(&pack).Error
}

func (p PackageRepository) UpdateTrafficUsage(ctx context.Context, req model.UpdateTrafficUsageRequest) error {
	err := p.db.
		WithContext(ctx).
		Model(&PackageEntity{}).
		Where("id = ?", req.ID).
		UpdateColumn("download_traffic_usage", gorm.Expr("download_traffic_usage + ?", req.DownloadTrafficUsage)).
		UpdateColumn("upload_traffic_usage", gorm.Expr("upload_traffic_usage + ?", req.UploadTrafficUsage)).Error
	if err != nil {
		return err
	}

	return nil
}

func (p PackageRepository) GetPackages(ctx context.Context, req model.GetPackagesRequest) (model.GetPackagesResponse, error) {
	var packages []PackageEntity

	query := p.db.WithContext(ctx).Model(&PackageEntity{})

	if req.UserID > 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	if err := query.Scopes(Paginate(&req.Pagination)).Order("created_at desc").Find(&packages).Error; err != nil {
		return model.GetPackagesResponse{}, err
	}

	return model.GetPackagesResponse{
		Packages:   toModelPackageEntities(packages),
		Pagination: req.Pagination,
	}, nil
}
