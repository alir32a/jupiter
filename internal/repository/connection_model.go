package repository

import (
	"github.com/alir32a/jupiter/internal/model"
	"time"
)

const (
	ConnectionsStatusConnected = "connected"
)

type ConnectionEntity struct {
	ID                   int
	Status               string
	Username             string
	ExternalID           string
	RemoteIP             string
	Location             string
	UserAgent            string
	Hostname             string
	DownloadTrafficUsage int
	UploadTrafficUsage   int
	ConnectedAt          time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func (ConnectionEntity) TableName() string {
	return "connection"
}

func toConnectionEntity(req model.ConnectionEntity) ConnectionEntity {
	return ConnectionEntity{
		ExternalID:           req.ExternalID,
		Username:             req.Username,
		RemoteIP:             req.RemoteIP,
		Location:             req.Location,
		UserAgent:            req.UserAgent,
		Hostname:             req.Hostname,
		DownloadTrafficUsage: req.DownloadTrafficUsage,
		UploadTrafficUsage:   req.UploadTrafficUsage,
		ConnectedAt:          req.ConnectedAt,
	}
}

func toConnectionEntities(req []model.ConnectionEntity) []ConnectionEntity {
	result := make([]ConnectionEntity, 0, len(req))

	for _, conn := range req {
		result = append(result, toConnectionEntity(conn))
	}

	return result
}

func toModelConnectionEntity(req ConnectionEntity) model.ConnectionEntity {
	return model.ConnectionEntity{
		ID:                   req.ID,
		ExternalID:           req.ExternalID,
		Username:             req.Username,
		Status:               req.Status,
		RemoteIP:             req.RemoteIP,
		Location:             req.Location,
		UserAgent:            req.UserAgent,
		Hostname:             req.Hostname,
		DownloadTrafficUsage: req.DownloadTrafficUsage,
		UploadTrafficUsage:   req.UploadTrafficUsage,
		ConnectedAt:          req.ConnectedAt,
		UpdatedAt:            req.UpdatedAt,
	}
}

func toModelConnectionEntities(req []ConnectionEntity) []model.ConnectionEntity {
	result := make([]model.ConnectionEntity, 0, len(req))

	for _, connection := range req {
		result = append(result, toModelConnectionEntity(connection))
	}

	return result
}
