package adapter

import (
	"context"
	"log/slog"

	"github.com/alnovi/sso/internal/entity"
)

type Notify interface {
}

type Repository interface {
	MigrateUp(ctx context.Context, log *slog.Logger) error
	MigrateDown(ctx context.Context, log *slog.Logger) error
	Close(ctx context.Context) error

	ClientByID(ctx context.Context, id string) (*entity.Client, error)

	UserByEmail(ctx context.Context, email string) (*entity.User, error)
}
