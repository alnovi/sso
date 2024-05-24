package migrations

import (
	"context"

	"github.com/alnovi/sso/internal/config"
)

const (
	DefaultClientAdminID   = "00000000-0000-0000-0000-000000000001"
	DefaultClientProfileID = "00000000-0000-0000-0000-000000000002"
	DefaultUserAdminID     = "00000000-0000-0000-0000-000000000001"
)

func EnvironmentIsTesting(ctx context.Context) bool {
	if env, ok := ctx.Value(config.KeyEnvironment).(string); ok {
		return env == config.EnvTesting
	}
	return false
}

func ClientAdminID(ctx context.Context) string {
	if adminID, ok := ctx.Value(config.KeyClientAdminID).(string); ok && adminID != "" {
		return adminID
	}
	return DefaultClientAdminID
}

func ClientProfileID(ctx context.Context) string {
	if profileID, ok := ctx.Value(config.KeyClientProfileID).(string); ok && profileID != "" {
		return profileID
	}
	return DefaultClientProfileID
}

func UserAdminID(ctx context.Context) string {
	if adminID, ok := ctx.Value(config.KeyUserAdminID).(string); ok && adminID != "" {
		return adminID
	}
	return DefaultUserAdminID
}
