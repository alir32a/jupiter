package service

import (
	"context"
	"errors"
	"github.com/alir32a/jupiter/config"
	"github.com/alir32a/jupiter/internal/errorext"
	"github.com/alir32a/jupiter/internal/model"
	"github.com/alir32a/jupiter/pkg/ocserv"
	"github.com/alir32a/jupiter/pkg/password"
	"github.com/alir32a/jupiter/pkg/util"
	clog "github.com/charmbracelet/log"
	"github.com/rs/xid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, req model.CreateUserRequest) (model.UserEntity, error)
	BanUser(ctx context.Context, id int) error
	UnbanUser(ctx context.Context, id int) error
	GetUsersByUsernames(ctx context.Context, usernames ...string) ([]model.UserEntity, error)
	GetUserByUsername(ctx context.Context, username string) (model.UserEntity, error)
	GetUsersStat(ctx context.Context) (model.GetUsersStatResponse, error)
	GetAllUsers(ctx context.Context, req model.GetAllUsersRequest) (model.GetAllUsersResponse, error)
	GetUserByID(ctx context.Context, id int) (model.UserEntity, error)
}

type UserPackageRepository interface {
	CreatePackage(ctx context.Context, req model.CreatePackageRequest) error
}

type UserService struct {
	cfg          *config.Config
	logger       *clog.Logger
	ocservClient *ocserv.Client
	repo         UserRepository
	packageRepo  UserPackageRepository
}

func NewUserService(cfg *config.Config, logger *clog.Logger, ocservClient *ocserv.Client, repo UserRepository,
	packageRepo UserPackageRepository) *UserService {
	return &UserService{
		cfg:          cfg,
		logger:       logger,
		ocservClient: ocservClient,
		repo:         repo,
		packageRepo:  packageRepo,
	}
}

func (u UserService) CreateUser(ctx context.Context, req model.CreateUserRequest) (model.CreateUserResponse, error) {
	req.ReferralCode = xid.New().String()

	user, err := u.repo.CreateUser(ctx, req)
	if err != nil {
		return model.CreateUserResponse{}, errorext.NewInternalError(u.logger, err)
	}

	pass := password.NewRandomPassword(password.DefaultLength)

	if err := u.ocservClient.CreateUser(ctx, req.Username, pass); err != nil {
		return model.CreateUserResponse{}, errorext.NewInternalError(u.logger, err)
	}

	if u.cfg.TrialPackage.Activated {
		err := u.packageRepo.CreatePackage(ctx, model.CreatePackageRequest{
			UserID:           user.ID,
			Traffic:          int(u.cfg.TrialPackage.TrafficLimit * util.GB),
			MaxConnections:   u.cfg.TrialPackage.MaxConnections,
			IsTrial:          true,
			ExpirationInDays: u.cfg.TrialPackage.ExpirationInDays,
		})
		if err != nil {
			u.logger.Error(err.Error())
		}
	}

	return model.CreateUserResponse{
		UserEntity: user,
		Password:   pass,
	}, nil
}

func (u UserService) ChangePassword(ctx context.Context, username string) (string, error) {
	_, err := u.repo.GetUsersByUsernames(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errorext.ErrUserNotFound
		}

		return "", errorext.NewInternalError(u.logger, err)
	}

	pass := password.NewRandomPassword(password.DefaultLength)
	err = u.ocservClient.ChangePassword(ctx, username, pass)
	if err != nil {
		return "", errorext.NewInternalError(u.logger, err)
	}

	return pass, nil
}

func (u UserService) BanUser(ctx context.Context, userID int) error {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := u.ocservClient.LockUser(ctx, user.Username); err != nil {
		return err
	}

	return u.repo.BanUser(ctx, userID)
}

func (u UserService) UnbanUser(ctx context.Context, userID int) error {
	user, err := u.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := u.ocservClient.UnlockUser(ctx, user.Username); err != nil {
		return err
	}

	return u.repo.UnbanUser(ctx, userID)
}

func (u UserService) GetUserByUsername(ctx context.Context, username string) (model.UserEntity, error) {
	return u.repo.GetUserByUsername(ctx, username)
}

func (u UserService) GetUsersStat(ctx context.Context) (model.GetUsersStatResponse, error) {
	return u.repo.GetUsersStat(ctx)
}

func (u UserService) GetAllUsers(ctx context.Context, req model.GetAllUsersRequest) (model.GetAllUsersResponse, error) {
	return u.repo.GetAllUsers(ctx, req)
}
