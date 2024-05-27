package usecase

import (
	"context"

	"github.com/alnovi/sso/internal/entity"
)

type repository interface {
	ClientByID(ctx context.Context, id string) (*entity.Client, error)
	ClientByIdAndSecret(ctx context.Context, id, secret string) (*entity.Client, error)
	UserByID(ctx context.Context, id string) (*entity.User, error)
	UserByEmail(ctx context.Context, email string) (*entity.User, error)
	//CreateToken(ctx context.Context, token *entity.Token) error
	ClientTokenByClassAndHash(ctx context.Context, clientID, class, hash string) (*entity.Token, error)
	DeleteTokenById(ctx context.Context, id string) error
}

type secure interface {
	NewCodeToken(ctx context.Context, user *entity.User, client *entity.Client) (*entity.Token, error)
	NewAccessToken(ctx context.Context, user *entity.User, client *entity.Client) (*entity.Token, error)
	NewRefreshToken(ctx context.Context, user *entity.User, client *entity.Client) (*entity.Token, error)
}

type UseCase struct {
	repo   repository
	secure secure
}

func New(repo repository, secure secure) *UseCase {
	return &UseCase{
		repo:   repo,
		secure: secure,
	}
}
