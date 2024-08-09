package repository

import (
	"context"
	"github.com/alir32a/jupiter/internal/model"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

func (a AdminRepository) CreateAdmin(ctx context.Context, req model.CreateAdminRequest) error {
	return a.db.WithContext(ctx).Model(&AdminEntity{}).Create(&req).Error
}

func (a AdminRepository) GetAdminByUsername(ctx context.Context, username string) (model.AdminEntity, error) {
	var admin AdminEntity

	err := a.db.WithContext(ctx).Model(&AdminEntity{}).First(&admin, "username = ?", username).Error
	if err != nil {
		return model.AdminEntity{}, err
	}

	return toModelAdminEntity(admin), nil
}

func (a AdminRepository) ChangePassword(ctx context.Context, username, password string) error {
	return a.db.
		WithContext(ctx).
		Model(&AdminEntity{}).
		Where("username = ?", username).
		UpdateColumn("password", password).Error
}
