package handler

import (
	"github.com/alir32a/jupiter/internal/model"
	"github.com/alir32a/jupiter/pkg/util"
	"time"
)

type ConnectionEntity struct {
	ID                   int       `json:"id"`
	ExternalID           string    `json:"external_id"`
	Username             string    `json:"username"`
	Status               string    `json:"status"`
	RemoteIP             string    `json:"remote_ip"`
	Location             string    `json:"location"`
	UserAgent            string    `json:"user_agent"`
	Hostname             string    `json:"hostname"`
	DownloadTrafficUsage int       `json:"download_traffic_usage"`
	UploadTrafficUsage   int       `json:"upload_traffic_usage"`
	ConnectedAt          time.Time `json:"connected_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type GetActiveConnectionsRequest struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Username string `query:"username"`
}

type GetActiveConnectionsResponse struct {
	Pagination
	Connections []ConnectionEntity `json:"connections"`
}

type GetSystemStatusResponse struct {
	TotalActiveConnections int    `json:"total_active_connections"`
	OnlineUsers            int    `json:"online_users"`
	TotalUsers             int    `json:"total_users"`
	TotalDownloadUsage     string `json:"total_download_usage"`
	TotalUploadUsage       string `json:"total_upload_usage"`
}

type DisconnectIDRequest struct {
	ID int `param:"id"`
}

type ConnectionSummary struct {
	UsedTraffic string    `json:"used_traffic"`
	IP          string    `json:"ip"`
	UserAgent   string    `json:"user_agent"`
	ConnectedAt time.Time `json:"connected_at"`
}

type GetUserActiveConnectionsRequest struct {
	Username string `query:"username"`
}

func toModelGetActiveConnectionsRequest(req GetActiveConnectionsRequest) model.GetActiveConnectionsRequest {
	return model.GetActiveConnectionsRequest{
		Pagination: &model.Pagination{
			CurrentPage: req.Page,
			PageSize:    req.PageSize,
		},
		Username: req.Username,
	}
}

func toCtrlConnectionEntity(req model.ConnectionEntity) ConnectionEntity {
	return ConnectionEntity{
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

func toCtrlConnectionEntities(connections []model.ConnectionEntity) []ConnectionEntity {
	result := make([]ConnectionEntity, 0, len(connections))

	for _, conn := range connections {
		result = append(result, toCtrlConnectionEntity(conn))
	}

	return result
}

func toCtrlGetSystemStatusResponse(req model.GetSystemStatusResponse) GetSystemStatusResponse {
	return GetSystemStatusResponse{
		TotalActiveConnections: req.TotalActiveConnections,
		OnlineUsers:            req.OnlineUsers,
		TotalUsers:             req.TotalUsers,
		TotalDownloadUsage:     util.ToHumanReadableBytes(req.TotalDownloadUsage),
		TotalUploadUsage:       util.ToHumanReadableBytes(req.TotalUploadUsage),
	}
}

func toCtrlConnectionSummary(conn model.ConnectionEntity) ConnectionSummary {
	return ConnectionSummary{
		UsedTraffic: util.ToHumanReadableBytes(conn.DownloadTrafficUsage + conn.UploadTrafficUsage),
		IP:          conn.RemoteIP,
		UserAgent:   conn.UserAgent,
		ConnectedAt: conn.ConnectedAt,
	}
}

func toCtrlConnectionSummaries(connections []model.ConnectionEntity) []ConnectionSummary {
	result := make([]ConnectionSummary, 0, len(connections))

	for _, conn := range connections {
		result = append(result, toCtrlConnectionSummary(conn))
	}

	return result
}
