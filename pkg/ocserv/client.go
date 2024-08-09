package ocserv

import (
	"context"
	"encoding/json"
	"os/exec"
)

type Client struct {
	passwordFilepath string
}

func (c Client) CreateUser(ctx context.Context, username, password string) error {
	cmd := exec.CommandContext(ctx, "ocpasswd", "-c", c.passwordFilepath, username)
	if cmd.Err != nil {
		return cmd.Err
	}

	_, err := cmd.Stdout.Write([]byte(password))
	if err != nil {
		return err
	}

	return nil
}

func (c Client) DisconnectUser(ctx context.Context, username string) error {
	return exec.CommandContext(ctx, "occtl", "disconnect", "user", username).Err
}

func (c Client) DisconnectID(ctx context.Context, id string) error {
	return exec.CommandContext(ctx, "occtl", "disconnect", "id", id).Err
}

func (c Client) GetConnections(ctx context.Context) ([]ConnectionEntity, error) {
	var connections []ConnectionEntity

	cmd := exec.CommandContext(ctx, "occtl", "-j", "show", "users")
	if cmd.Err != nil {
		return nil, cmd.Err
	}

	data, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &connections); err != nil {
		return nil, err
	}

	return connections, nil
}

func (c Client) ShutdownServer(ctx context.Context) error {
	return exec.CommandContext(ctx, "occtl", "stop", "now").Err
}

func (c Client) ChangePassword(ctx context.Context, username, password string) error {
	err := exec.CommandContext(ctx, "ocpasswd", "-c", c.passwordFilepath, "--delete", username).Err
	if err != nil {
		return err
	}

	return c.CreateUser(ctx, username, password)
}

func (c Client) LockUser(ctx context.Context, username string) error {
	return exec.CommandContext(ctx, "ocpasswd", "-c", c.passwordFilepath, "-l", username).Err
}

func (c Client) UnlockUser(ctx context.Context, username string) error {
	return exec.CommandContext(ctx, "ocpasswd", "-c", c.passwordFilepath, "-u", username).Err
}

func NewClient(passwordFilepath string) *Client {
	return &Client{passwordFilepath: passwordFilepath}
}

func CheckInstallation(ctx context.Context) error {
	err := exec.CommandContext(ctx, "ocserv").Err
	if err != nil {
		return err
	}

	err = exec.CommandContext(ctx, "occtl").Err
	if err != nil {
		return err
	}

	return exec.CommandContext(ctx, "ocpasswd").Err
}
