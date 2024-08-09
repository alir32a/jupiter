package service

import (
	"cmp"
	"context"
	"github.com/alir32a/jupiter/internal/model"
	"github.com/alir32a/jupiter/internal/util"
	"github.com/alir32a/jupiter/pkg/ocserv"
	clog "github.com/charmbracelet/log"
	"math"
	"slices"
	"time"
)

type ConnectionRepository interface {
	UpsertConnections(ctx context.Context, req model.UpsertConnectionsRequest) error
	GetActiveConnections(ctx context.Context, req model.GetActiveConnectionsRequest) (model.GetActiveConnectionsResponse, error)
	Disconnect(ctx context.Context, req model.DisconnectRequest) error
	GetUserActiveConnections(ctx context.Context, username string) ([]model.ConnectionEntity, error)
	GetSystemStatus(ctx context.Context) (model.GetSystemStatusResponse, error)
	DisconnectID(ctx context.Context, id int) (string, error)
}

type ConnectionPackageRepository interface {
	GetUsersActiveAndReservedPackages(ctx context.Context, userIDs ...int) (model.GetUsersActivePackagesResponse, error)
	UpdateTrafficUsage(ctx context.Context, req model.UpdateTrafficUsageRequest) error
}

type ConnectionUserRepository interface {
	GetUsersByUsernames(ctx context.Context, usernames ...string) ([]model.UserEntity, error)
	GetTotalUsersCount(ctx context.Context) (int, error)
}

type ConnectionService struct {
	logger       *clog.Logger
	ocservClient *ocserv.Client
	repo         ConnectionRepository
	packageRepo  ConnectionPackageRepository
	userRepo     ConnectionUserRepository
}

func NewConnectionService(logger *clog.Logger, ocservClient *ocserv.Client, repo ConnectionRepository,
	packageRepo ConnectionPackageRepository, userRepo ConnectionUserRepository) *ConnectionService {
	return &ConnectionService{
		logger:       logger,
		ocservClient: ocservClient,
		repo:         repo,
		packageRepo:  packageRepo,
		userRepo:     userRepo,
	}
}

func (c ConnectionService) UpsertConnections(ctx context.Context, req model.UpsertConnectionsRequest) error {
	return c.repo.UpsertConnections(ctx, req)
}

func (c ConnectionService) ManageActiveConnections(ctx context.Context, lastUpdated time.Time) error {
	resp, err := c.repo.GetActiveConnections(ctx, model.GetActiveConnectionsRequest{})
	if err != nil {
		return err
	}

	connections := resp.Connections

	connectionsByUsers := util.MapElementsBy(connections, func(conn model.ConnectionEntity) string {
		return conn.Username
	})

	users, err := c.userRepo.GetUsersByUsernames(ctx, util.MapKeys(connectionsByUsers)...)
	if err != nil {
		return err
	}

	usersMap := util.MapUniqueElementsBy(users, func(user model.UserEntity) string {
		return user.Username
	})

	packages, err := c.packageRepo.GetUsersActiveAndReservedPackages(ctx, util.SliceMap(users, func(user model.UserEntity) int {
		return user.ID
	})...)
	if err != nil {
		return err
	}

	packagesMap := util.MapUniqueElementsBy(packages.Packages, func(userPacks model.GetUserPackages) int {
		return userPacks.UserID
	})

	for username, connections := range connectionsByUsers {
		user, ok := usersMap[username]
		if !ok {
			c.DisconnectUser(ctx, username, connections)

			continue
		}

		if user.BannedAt != nil {
			c.DisconnectUser(ctx, username, connections)

			continue
		}

		packs, ok := packagesMap[user.ID]
		if !ok {
			c.DisconnectUser(ctx, username, connections)

			continue
		}

		totalUsage := util.SumFunc(connections, func(conn model.ConnectionEntity) int {
			return conn.DownloadTrafficUsage + conn.UploadTrafficUsage
		})

		totalDownloadUsage := util.SumFunc(connections, func(conn model.ConnectionEntity) int {
			return conn.DownloadTrafficUsage
		})

		totalUploadUsage := util.SumFunc(connections, func(conn model.ConnectionEntity) int {
			return conn.DownloadTrafficUsage
		})

		remainingTraffic := getRemainingTraffic(packs.ActivePackage, totalUsage)
		if remainingTraffic < 0 {
			if packs.ActivePackage.ID != 0 {
				c.DisconnectUser(ctx, username, connections)

				continue
			}

			totalUsageDiff := totalUsage + remainingTraffic

			downloadUsage, uploadUsage := calculateDownloadAndUploadByTotalUsage(totalDownloadUsage, totalUploadUsage, totalUsageDiff)
			err = c.packageRepo.UpdateTrafficUsage(ctx, model.UpdateTrafficUsageRequest{
				ID:                   packs.ActivePackage.ID,
				DownloadTrafficUsage: downloadUsage,
				UploadTrafficUsage:   uploadUsage,
			})
			if err != nil {
				return err
			}

			totalDownloadUsage -= downloadUsage
			totalUploadUsage -= uploadUsage

			if len(packs.ReservedPackages) <= 0 {
				c.DisconnectUser(ctx, username, connections)

				continue
			}

			slices.SortFunc(packs.ReservedPackages, func(a, b model.PackageEntity) int {
				return cmp.Compare(a.ExpirationInDays, b.ExpirationInDays)
			})

			remainingTraffic = int(math.Abs(float64(remainingTraffic)))
			for _, pack := range packs.ReservedPackages {
				if pack.TrafficLimit > remainingTraffic {
					downloadUsage, uploadUsage := calculateDownloadAndUploadByTotalUsage(totalDownloadUsage, totalUploadUsage, remainingTraffic)
					err = c.packageRepo.UpdateTrafficUsage(ctx, model.UpdateTrafficUsageRequest{
						ID:                   pack.ID,
						DownloadTrafficUsage: downloadUsage,
						UploadTrafficUsage:   uploadUsage,
					})
					if err != nil {
						return err
					}

					remainingTraffic = 0
					break
				}

				downloadUsage, uploadUsage := calculateDownloadAndUploadByTotalUsage(totalDownloadUsage, totalUploadUsage, pack.TrafficLimit)
				err = c.packageRepo.UpdateTrafficUsage(ctx, model.UpdateTrafficUsageRequest{
					ID:                   pack.ID,
					DownloadTrafficUsage: downloadUsage,
					UploadTrafficUsage:   uploadUsage,
				})
				if err != nil {
					return err
				}

				remainingTraffic -= downloadUsage + uploadUsage
			}

			if remainingTraffic > 0 {
				c.DisconnectUser(ctx, username, connections)

				continue
			}
		}

		for _, conn := range connections {
			if conn.UpdatedAt.Before(lastUpdated) {
				err := c.repo.Disconnect(ctx, model.DisconnectRequest{
					ConnectionID:         conn.ID,
					DownloadTrafficUsage: conn.DownloadTrafficUsage,
					UploadTrafficUsage:   conn.UploadTrafficUsage,
					Username:             conn.Username,
				})
				if err != nil {
					c.logger.Error(err)
				}
			}
		}
	}

	return nil
}

func (c ConnectionService) DisconnectUser(ctx context.Context, username string, connections []model.ConnectionEntity) {
	if err := c.ocservClient.DisconnectUser(ctx, username); err != nil {
		c.logger.Error(err)

		return
	}

	if err := c.ocservClient.LockUser(ctx, username); err != nil {
		c.logger.Error(err)

		return
	}

	for _, conn := range connections {
		err := c.repo.Disconnect(ctx, model.DisconnectRequest{
			ConnectionID:         conn.ID,
			DownloadTrafficUsage: conn.DownloadTrafficUsage,
			UploadTrafficUsage:   conn.UploadTrafficUsage,
			Username:             conn.Username,
		})
		if err != nil {
			c.logger.Error(err)
		}
	}
}

func (c ConnectionService) DisconnectID(ctx context.Context, id int) error {
	externalID, err := c.repo.DisconnectID(ctx, id)
	if err != nil {
		return err
	}

	return c.ocservClient.DisconnectID(ctx, externalID)
}

func (c ConnectionService) GetActiveConnections(ctx context.Context, req model.GetActiveConnectionsRequest) (model.GetActiveConnectionsResponse, error) {
	return c.repo.GetActiveConnections(ctx, req)
}

func (c ConnectionService) GetUserActiveConnections(ctx context.Context, username string) ([]model.ConnectionEntity, error) {
	return c.repo.GetUserActiveConnections(ctx, username)
}

func (c ConnectionService) GetSystemStatus(ctx context.Context) (model.GetSystemStatusResponse, error) {
	resp, err := c.repo.GetSystemStatus(ctx)
	if err != nil {
		return model.GetSystemStatusResponse{}, err
	}

	totalUsers, err := c.userRepo.GetTotalUsersCount(ctx)
	if err != nil {
		return model.GetSystemStatusResponse{}, err
	}
	resp.TotalUsers = totalUsers

	return resp, nil
}

func getRemainingTraffic(pack model.PackageEntity, usage int) int {
	return pack.TrafficLimit - (pack.DownloadTrafficUsage + pack.UploadTrafficUsage + usage)
}

func calculateDownloadAndUploadByTotalUsage(totalDownloadUsage, totalUploadUsage, usage int) (int, int) {
	totalUsage := float64(totalDownloadUsage + totalUploadUsage)
	downloadShare := float64(totalDownloadUsage) / totalUsage
	uploadShare := float64(totalUploadUsage) / totalUsage

	return int(downloadShare * float64(usage)), int(uploadShare * float64(usage))
}
