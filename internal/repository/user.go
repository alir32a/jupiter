package repository

import (
	"context"
	"errors"
	"github.com/alir32a/jupiter/internal/errorext"
	"github.com/alir32a/jupiter/internal/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u UserRepository) CreateUser(ctx context.Context, req model.CreateUserRequest) (model.UserEntity, error) {
	user := UserEntity{
		Username:     req.Username,
		ExternalID:   req.ExternalID,
		UserType:     req.UserType,
		ReferralCode: req.ReferralCode,
		Referral:     req.Referral,
	}

	err := u.db.
		WithContext(ctx).
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "id"}}}).
		Create(&user).Error
	if err != nil {
		return model.UserEntity{}, err
	}

	return toModelUserEntity(user), nil
}

func (u UserRepository) GetUsersByUsernames(ctx context.Context, usernames ...string) ([]model.UserEntity, error) {
	var users []UserEntity

	err := u.db.WithContext(ctx).Where("username in ?", usernames).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return toModelUserEntities(users), nil
}

func (u UserRepository) GetUserByID(ctx context.Context, id int) (model.UserEntity, error) {
	var user UserEntity

	err := u.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.UserEntity{}, errorext.NewNotFoundError(errorext.ErrUserNotFound)
		}

		return model.UserEntity{}, err
	}

	return toModelUserEntity(user), nil
}

func (u UserRepository) GetUsersByIDs(ctx context.Context, ids []int) ([]model.UserEntity, error) {
	var users []UserEntity

	err := u.db.WithContext(ctx).Where("id in ?", ids).Find(&users).Error
	if err != nil {
		return nil, err
	}

	return toModelUserEntities(users), nil
}

func (u UserRepository) GetUserByUsername(ctx context.Context, username string) (model.UserEntity, error) {
	var user UserEntity

	err := u.db.WithContext(ctx).First(&user, "username = ?", username).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.UserEntity{}, errorext.NewNotFoundError(errorext.ErrUserNotFound)
		}

		return model.UserEntity{}, err
	}

	return toModelUserEntity(user), nil
}

func (u UserRepository) BanUser(ctx context.Context, id int) error {
	return u.db.WithContext(ctx).Model(&UserEntity{}).Where("id = ?", id).UpdateColumn("banned_at", time.Now()).Error
}

func (u UserRepository) UnbanUser(ctx context.Context, id int) error {
	return u.db.WithContext(ctx).Model(&UserEntity{}).Where("id = ?", id).UpdateColumn("banned_at", nil).Error
}

func (u UserRepository) GetUsersStat(ctx context.Context) (model.GetUsersStatResponse, error) {
	var result = model.GetUsersStatResponse{}

	query := u.db.WithContext(ctx).Model(&UserEntity{})

	if err := query.Select("count(*)").Scan(&result.TotalUsers).Error; err != nil {
		return model.GetUsersStatResponse{}, err
	}

	if err := query.Select("count(*)").Where("banned_at is not null").Scan(&result.TotalBannedUsers).Error; err != nil {
		return model.GetUsersStatResponse{}, err
	}

	err := u.db.WithContext(ctx).Model(&UserEntity{}).
		Select("count(distinct \"user\".id)").
		Joins("inner join package on package.user_id = \"user\".id").
		Where(`package.download_traffic_usage + package.upload_traffic_usage < package.traffic_limit and 
        (package.expire_at > now() or package.expire_at is null)`).
		Scan(&result.TotalActiveUsers).Error
	if err != nil {
		return model.GetUsersStatResponse{}, err
	}

	return result, nil
}

func (u UserRepository) GetAllUsers(ctx context.Context, req model.GetAllUsersRequest) (model.GetAllUsersResponse, error) {
	var users []UserEntity

	query := u.db.WithContext(ctx).Model(&UserEntity{})

	if req.Username != "" {
		query = query.Where("username = ?", req.Username)
	}

	err := query.Scopes(Paginate(&req.Pagination)).Order("created_at desc").Find(&users).Error
	if err != nil {
		return model.GetAllUsersResponse{}, err
	}

	return model.GetAllUsersResponse{
		Users:      toModelUserEntities(users),
		Pagination: req.Pagination,
	}, nil
}

func (u UserRepository) GetTotalUsersCount(ctx context.Context) (int, error) {
	var total int64

	if err := u.db.WithContext(ctx).Model(&UserEntity{}).Count(&total).Error; err != nil {
		return 0, err
	}

	return int(total), nil
}
