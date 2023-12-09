package repository

import (
	"context"

	"github.com/alnovi/sso/internal/entity"
)

type Repository interface {
	MigrateUp() error
	Close() error

	Begin(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error

	GetUserById(ctx context.Context, id string) (*entity.User, error)
	GetUserByLoginOrEmail(ctx context.Context, login string) (*entity.User, error)

	GetClientByID(ctx context.Context, id string) (*entity.Client, error)
	GetClientByClass(ctx context.Context, class string) (*entity.Client, error)

	CanUseClient(ctx context.Context, client entity.Client, user entity.User) error

	CreateToken(ctx context.Context, token *entity.Token) error
	UpdateToken(ctx context.Context, token *entity.Token) error
	GetTokenByHash(ctx context.Context, hash string) (*entity.Token, error)
}
