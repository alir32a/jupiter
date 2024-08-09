package repository

import (
	"github.com/alir32a/jupiter/internal/model"
	"time"
)

type AdminEntity struct {
	ID        int
	Username  string
	Password  string
	IsActive  bool
	CreatedAt time.Time
}

func (AdminEntity) TableName() string {
	return "admin"
}

func toModelAdminEntity(req AdminEntity) model.AdminEntity {
	return model.AdminEntity{
		ID:        req.ID,
		Username:  req.Username,
		Password:  req.Password,
		IsActive:  req.IsActive,
		CreatedAt: req.CreatedAt,
	}
}
