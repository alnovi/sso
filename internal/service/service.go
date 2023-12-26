package service

import (
	"context"

	"github.com/alnovi/sso/internal/entity"
)

type Auth interface {
	AuthByCredentials(ctx context.Context, login, password string) (*entity.User, error)
}

type Client interface {
	GetManagerClient(ctx context.Context) (*entity.Client, error)
	GetProfileClient(ctx context.Context) (*entity.Client, error)
	GetClient(ctx context.Context, id string, callback, secret *string) (*entity.Client, error)
}

type Token interface {
	NewCode(ctx context.Context, user entity.User, client entity.Client, meta *entity.TokenMeta) (*entity.Token, error)
	NewAccess(ctx context.Context, user entity.User, client entity.Client) (*entity.Token, error)
	NewRefresh(ctx context.Context, user entity.User, client entity.Client, meta *entity.TokenMeta) (*entity.Token, error)
	Validate(ctx context.Context, clientId, secret, class, hash string) (*entity.Token, error)
	FindToken(ctx context.Context, client entity.Client, class, hash string) (*entity.Token, error)
	RemoveToken(ctx context.Context, tokenId string) error
}

type User interface {
	CanUseClient(ctx context.Context, user entity.User, client entity.Client) error
	GetUserById(ctx context.Context, userId string) (*entity.User, error)
}
