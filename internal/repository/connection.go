package repository

import (
	"context"
	"github.com/alir32a/jupiter/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type ConnectionRepository struct {
	db *gorm.DB
}

func NewConnectionRepository(db *gorm.DB) *ConnectionRepository {
	return &ConnectionRepository{db: db}
}

func (c ConnectionRepository) UpsertConnections(ctx context.Context, req model.UpsertConnectionsRequest) error {
	connections := toConnectionEntities(req.Connections)

	return c.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"download_traffic_usage", "upload_traffic_usage", "updated_at"}),
	}).Create(&connections).Error
}

func (c ConnectionRepository) GetActiveConnections(ctx context.Context, req model.GetActiveConnectionsRequest) (model.GetActiveConnectionsResponse, error) {
	var result []ConnectionEntity

	query := c.db.
		WithContext(ctx).
		Model(&ConnectionEntity{}).
		Where("status = ?", ConnectionsStatusConnected)

	if req.Pagination != nil {
		query = query.Scopes(Paginate(req.Pagination))
	}

	if req.Username != "" {
		query = query.Where("username = ?", req.Username)
	}

	err := query.Order("connected_at desc").Find(&result).Error
	if err != nil {
		return model.GetActiveConnectionsResponse{}, err
	}

	return model.GetActiveConnectionsResponse{
		Connections: toModelConnectionEntities(result),
		Pagination:  req.Pagination,
	}, nil
}

func (c ConnectionRepository) GetUserActiveConnections(ctx context.Context, username string) ([]model.ConnectionEntity, error) {
	var result []ConnectionEntity

	err := c.db.
		WithContext(ctx).
		Model(&ConnectionEntity{}).
		Where("status = ?", ConnectionsStatusConnected).
		Find(&result, "username = ?", username).Error
	if err != nil {
		return nil, err
	}

	return toModelConnectionEntities(result), nil
}

func (c ConnectionRepository) Disconnect(ctx context.Context, req model.DisconnectRequest) error {
	return c.db.
		WithContext(ctx).
		Model(&ConnectionEntity{}).
		Where("id = ?", req.ConnectionID).
		Updates(&ConnectionEntity{
			DownloadTrafficUsage: req.DownloadTrafficUsage,
			UploadTrafficUsage:   req.UploadTrafficUsage,
			UpdatedAt:            time.Now(),
			Status:               model.ConnectionStatusDisConnected,
		}).Error
}

func (c ConnectionRepository) DisconnectID(ctx context.Context, id int) (string, error) {
	req := ConnectionEntity{Status: model.ConnectionStatusDisConnected}

	err := c.db.
		WithContext(ctx).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "external_id"}}}).
		Where("id = ?", id).
		Updates(&req).Error
	if err != nil {
		return "", err
	}

	return req.ExternalID, nil
}

func (c ConnectionRepository) GetSystemStatus(ctx context.Context) (model.GetSystemStatusResponse, error) {
	var result model.GetSystemStatusResponse

	err := c.db.
		WithContext(ctx).
		Model(&ConnectionEntity{}).
		Select("count(*) as total_active_connections, count(distinct username) as online_users").
		Where("status = ?", ConnectionsStatusConnected).
		Find(&result).Error
	if err != nil {
		return model.GetSystemStatusResponse{}, err
	}

	err = c.db.
		WithContext(ctx).
		Model(&ConnectionEntity{}).
		Select("sum(download_traffic_usage) as total_download_usage, sum(upload_traffic_usage) as total_upload_usage").
		Find(&result).Error
	if err != nil {
		return model.GetSystemStatusResponse{}, err
	}

	return result, nil
}
