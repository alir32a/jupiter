package handler

import "github.com/alir32a/jupiter/internal/model"

const AccessTokenCookieName = "adminAccess"

type AdminLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SelfResponse struct {
	Username string `json:"username"`
}

type ChangePasswordRequest struct {
	Username        string `json:"-"`
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}

func toModelAdminLoginRequest(req AdminLoginRequest) model.AdminLoginRequest {
	return model.AdminLoginRequest{
		Username: req.Username,
		Password: req.Password,
	}
}

func toModelChangePasswordRequest(req ChangePasswordRequest) model.ChangePasswordRequest {
	return model.ChangePasswordRequest{
		Username:        req.Username,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}
}
