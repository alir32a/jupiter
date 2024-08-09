package model

import "time"

type AdminEntity struct {
	ID        int
	Username  string
	Password  string
	IsActive  bool
	CreatedAt time.Time
}

type CreateAdminRequest struct {
	Username string
	Password string
}

type AdminLoginRequest struct {
	Username string
	Password string
}

type ChangePasswordRequest struct {
	Username        string
	CurrentPassword string
	NewPassword     string
}
