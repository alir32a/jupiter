package model

import "time"

const (
	UserTypeTelegram = "telegram"
)

type CreateUserRequest struct {
	Username     string
	ExternalID   string
	UserType     string
	ReferralCode string
	Referral     *string
}

type CreateUserResponse struct {
	UserEntity
	Password string
}

type GetUserRequest struct {
	ID         int
	Username   string
	ExternalID string
}

type UserEntity struct {
	ID           int
	Username     string
	ExternalID   string
	UserType     string
	ReferralCode string
	Referral     *string
	BannedAt     *time.Time
	CreatedAt    time.Time
}

type GetUsersStatResponse struct {
	TotalUsers       int
	TotalActiveUsers int
	TotalBannedUsers int
}

type GetAllUsersRequest struct {
	Pagination
	Username string
}

type GetAllUsersResponse struct {
	Users []UserEntity
	Pagination
}
