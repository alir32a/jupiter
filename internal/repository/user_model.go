package repository

import (
	"github.com/alir32a/jupiter/internal/model"
	"time"
)

type UserEntity struct {
	ID           int
	Username     string
	ExternalID   string
	UserType     string
	ReferralCode string
	Referral     *string
	BannedAt     *time.Time
	CreatedAt    time.Time
	DeletedAt    *time.Time
}

func (UserEntity) TableName() string {
	return "user"
}

func toModelUserEntity(req UserEntity) model.UserEntity {
	return model.UserEntity{
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

func toModelUserEntities(req []UserEntity) []model.UserEntity {
	result := make([]model.UserEntity, 0, len(req))

	for _, user := range req {
		result = append(result, toModelUserEntity(user))
	}

	return result
}
