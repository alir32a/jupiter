package service

import (
	"context"
	"errors"
	"github.com/alir32a/jupiter/internal/errorext"
	"github.com/alir32a/jupiter/internal/model"
	"github.com/alir32a/jupiter/pkg/password"
	clog "github.com/charmbracelet/log"
)

type AdminRepository interface {
	CreateAdmin(ctx context.Context, req model.CreateAdminRequest) error
	GetAdminByUsername(ctx context.Context, username string) (model.AdminEntity, error)
	ChangePassword(ctx context.Context, username, password string) error
}

type AdminService struct {
	repo   AdminRepository
	logger *clog.Logger
}

func NewAdminService(repo AdminRepository, logger *clog.Logger) *AdminService {
	return &AdminService{
		repo:   repo,
		logger: logger,
	}
}

func (a AdminService) CreateAdmin(ctx context.Context, req model.CreateAdminRequest) error {
	var err error

	req.Password, err = password.HashPassword(req.Password)
	if err != nil {
		return err
	}

	return a.repo.CreateAdmin(ctx, req)
}

func (a AdminService) Login(ctx context.Context, req model.AdminLoginRequest) error {
	admin, err := a.repo.GetAdminByUsername(ctx, req.Username)
	if err != nil {
		return errorext.NewNotFoundError(errorext.ErrUserOrPasswordIsIncorrect)
	}

	if err := password.ComparePasswords(admin.Password, req.Password); err != nil {
		return errorext.NewNotFoundError(errorext.ErrUserOrPasswordIsIncorrect)
	}

	return nil
}

func (a AdminService) ChangePassword(ctx context.Context, req model.ChangePasswordRequest) error {
	user, err := a.repo.GetAdminByUsername(ctx, req.Username)
	if err != nil {
		return err
	}

	if err := password.ComparePasswords(user.Password, req.CurrentPassword); err != nil {
		return errorext.NewBadRequestError(errors.New("current password is wrong"))
	}

	pass, err := password.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return a.repo.ChangePassword(ctx, req.Username, pass)
}
