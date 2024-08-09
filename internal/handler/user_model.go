package handler

import (
	"github.com/alir32a/jupiter/internal/model"
	"time"
)

type UserEntity struct {
	ID           int        `json:"id"`
	Username     string     `json:"username"`
	ExternalID   string     `json:"external_id"`
	UserType     string     `json:"user_type"`
	ReferralCode string     `json:"referral_code"`
	Referral     *string    `json:"referral"`
	BannedAt     *time.Time `json:"banned_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type GetAllUsersRequest struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Username string `query:"username"`
}

type GetAllUsersResponse struct {
	Users []UserEntity `json:"users"`
	Pagination
}

type BanUserRequest struct {
	ID int `param:"id"`
}

type UnBanUserRequest struct {
	ID int `param:"id"`
}

func toModelGetAllUsersRequest(req GetAllUsersRequest) model.GetAllUsersRequest {
	return model.GetAllUsersRequest{
		Pagination: model.Pagination{
			CurrentPage: req.Page,
			PageSize:    req.PageSize,
		},
		Username: req.Username,
	}
}

func toCtrlUserEntity(req model.UserEntity) UserEntity {
	return UserEntity{
		ID:           req.ID,
		Username:     req.Username,
		ExternalID:   req.ExternalID,
		UserType:     req.UserType,
		ReferralCode: req.ReferralCode,
		Referral:     req.Referral,
		BannedAt:     req.BannedAt,
		CreatedAt:    req.CreatedAt,
	}
}

func toCtrlUserEntities(users []model.UserEntity) []UserEntity {
	result := make([]UserEntity, 0, len(users))

	for _, user := range users {
		result = append(result, toCtrlUserEntity(user))
	}

	return result
}

func toCtrlGetAllUsersResponse(req model.GetAllUsersResponse) GetAllUsersResponse {
	return GetAllUsersResponse{
		Users:      toCtrlUserEntities(req.Users),
		Pagination: toCtrlPagination(req.Pagination),
	}
}
