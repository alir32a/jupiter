package model

import "time"

type CreatePackageRequest struct {
	Username         string
	UserID           int
	Traffic          int
	MaxConnections   int
	IsTrial          bool
	ExpirationInDays int
	ExpireAt         *time.Time
}

type PackageEntity struct {
	ID                   int
	UserID               int
	Username             string
	TrafficLimit         int
	DownloadTrafficUsage int
	UploadTrafficUsage   int
	MaxConnections       int
	IsTrial              bool
	ExpirationInDays     int
	ExpireAt             *time.Time
	CreatedAt            time.Time
}

type UpdateTrafficUsageRequest struct {
	ID                   int
	DownloadTrafficUsage int
	UploadTrafficUsage   int
}

type GetUserPackages struct {
	UserID           int
	ActivePackage    PackageEntity
	ReservedPackages []PackageEntity
}

type GetUsersActivePackagesResponse struct {
	Packages []GetUserPackages
}

type GetPackagesRequest struct {
	Pagination
	Username string
	UserID   int
}

type GetPackagesResponse struct {
	Packages []PackageEntity
	Pagination
}
