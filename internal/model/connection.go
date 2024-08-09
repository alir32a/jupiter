package model

import "time"

const (
	ConnectionStatusConnected    = "connected"
	ConnectionStatusDisConnected = "disconnected"
)

type ConnectionEntity struct {
	ID                   int
	ExternalID           string
	Username             string
	Status               string
	RemoteIP             string
	Location             string
	UserAgent            string
	Hostname             string
	DownloadTrafficUsage int
	UploadTrafficUsage   int
	ConnectedAt          time.Time
	UpdatedAt            time.Time
}

type UpsertConnectionsRequest struct {
	Connections []ConnectionEntity
}

type DisconnectRequest struct {
	ConnectionID         int
	DownloadTrafficUsage int
	UploadTrafficUsage   int
	Username             string
}

type GetActiveConnectionsRequest struct {
	*Pagination
	Username string
}

type GetActiveConnectionsResponse struct {
	Connections []ConnectionEntity
	*Pagination
}

type GetSystemStatusResponse struct {
	TotalActiveConnections int
	OnlineUsers            int
	TotalUsers             int
	TotalDownloadUsage     int
	TotalUploadUsage       int
}
