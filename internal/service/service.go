package service

import (
	"context"

	"github.com/alnovi/sso/internal/entity"
)

type Auth interface {
	AuthByCredentials(ctx context.Context, login, password string) (*entity.User, error)
	CanUseClient(ctx context.Context, user entity.User, client entity.Client) error
}

type Client interface {
	GetProfileClient(ctx context.Context) (*entity.Client, error)
	GetClient(ctx context.Context, id string, callback, secret *string) (*entity.Client, error)
}

type Token interface {
	NewCode(ctx context.Context, userId, clientId, ip, agent string) (*entity.Token, error)
}

type User interface {
	GetUserById(ctx context.Context, userId string) (*entity.User, error)
}
