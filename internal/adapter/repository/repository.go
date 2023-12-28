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
	DeleteToken(ctx context.Context, tokenId string) error
	GetTokenByClassAndHash(ctx context.Context, class, hash string) (*entity.Token, error)
	GetTokenByClientAndHash(ctx context.Context, clientId, hash string) (*entity.Token, error)
	TokensByUser(ctx context.Context, userId string, class *string) ([]*entity.Token, error)
}
