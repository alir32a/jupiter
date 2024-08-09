package service

import (
	"context"
	"errors"
	"github.com/alir32a/jupiter/internal/errorext"
	"github.com/alir32a/jupiter/internal/model"
	"github.com/alir32a/jupiter/pkg/util"
	clog "github.com/charmbracelet/log"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

type PackageUserRepository interface {
	GetUsersByUsernames(ctx context.Context, usernames ...string) ([]model.UserEntity, error)
	GetUserByUsername(ctx context.Context, username string) (model.UserEntity, error)
	GetUsersByIDs(ctx context.Context, ids []int) ([]model.UserEntity, error)
}

type PackageRepository interface {
	GetUserActiveAndReservedPackages(ctx context.Context, userID int) (model.GetUserPackages, error)
	GetPackages(ctx context.Context, req model.GetPackagesRequest) (model.GetPackagesResponse, error)
	CreatePackage(ctx context.Context, req model.CreatePackageRequest) error
}

type PackageService struct {
	repo     PackageRepository
	userRepo PackageUserRepository
	logger   *clog.Logger
}

func NewPackageService(logger *clog.Logger, repo PackageRepository, userRepo PackageUserRepository) *PackageService {
	return &PackageService{
		logger:   logger,
		repo:     repo,
		userRepo: userRepo,
	}
}

func (p PackageService) GetUserActivePackages(ctx context.Context, username string) (model.GetUserPackages, error) {
	resp, err := p.userRepo.GetUsersByUsernames(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.GetUserPackages{}, errorext.ErrUserNotFound
		}

		return model.GetUserPackages{}, errorext.NewInternalError(p.logger, err)
	}

	if len(resp) <= 0 {
		return model.GetUserPackages{}, errorext.ErrUserNotFound
	}

	user := resp[0]
	if user.BannedAt != nil {
		return model.GetUserPackages{}, errorext.ErrUserBanned
	}

	packages, err := p.repo.GetUserActiveAndReservedPackages(ctx, user.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.GetUserPackages{}, errorext.ErrNoActivePackage
		}

		return model.GetUserPackages{}, errorext.NewInternalError(p.logger, err)
	}

	return packages, nil
}

func (p PackageService) GetPackages(ctx context.Context, req model.GetPackagesRequest) (model.GetPackagesResponse, error) {
	if req.Username != "" {
		user, err := p.userRepo.GetUserByUsername(ctx, req.Username)
		if err != nil {
			if errors.Is(err, errorext.ErrUserNotFound) {
				return model.GetPackagesResponse{}, nil
			}

			return model.GetPackagesResponse{}, err
		}

		req.UserID = user.ID
	}

	packages, err := p.repo.GetPackages(ctx, req)
	if err != nil {
		return model.GetPackagesResponse{}, err
	}

	packages.Packages, err = p.fillPackagesUsernames(ctx, packages.Packages)
	if err != nil {
		return model.GetPackagesResponse{}, err
	}

	return packages, nil
}

func (p PackageService) fillPackagesUsernames(ctx context.Context, packages []model.PackageEntity) ([]model.PackageEntity, error) {
	result := make([]model.PackageEntity, 0, len(packages))

	mapPackagesByUserID := util.MapStructsByField(packages, func(pack model.PackageEntity) int {
		return pack.UserID
	})

	users, err := p.userRepo.GetUsersByIDs(ctx, maps.Keys(mapPackagesByUserID))
	if err != nil {
		return nil, err
	}

	usersMap := util.MapStructsByUniqueField(users, func(user model.UserEntity) int {
		return user.ID
	})

	for userID, packs := range mapPackagesByUserID {
		user, ok := usersMap[userID]
		if !ok {
			continue
		}

		util.MapSlice(packs, func(pack model.PackageEntity) model.PackageEntity {
			pack.Username = user.Username

			return pack
		})

		result = append(result, packs...)
	}

	return result, nil
}

func (p PackageService) CreatePackage(ctx context.Context, req model.CreatePackageRequest) error {
	user, err := p.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return err
	}

	req.UserID = user.ID
	req.Traffic *= util.GB

	return p.repo.CreatePackage(ctx, req)
}

func (p PackageService) GetUserActiveAndReservedPackages(ctx context.Context, userID int) (model.GetUserPackages, error) {
	return p.repo.GetUserActiveAndReservedPackages(ctx, userID)
}
