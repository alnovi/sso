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
	UpdateUser(ctx context.Context, user *entity.User) error
	TokenByClassAndHash(ctx context.Context, class, hash string) (*entity.Token, error)
	ClientTokenByClassAndHash(ctx context.Context, clientID, class, hash string) (*entity.Token, error)
	DeleteTokenById(ctx context.Context, id string) error
}

type notify interface {
	ResetPassword(ctx context.Context, user *entity.User, token *entity.Token) error
}

type secure interface {
	NewCodeToken(ctx context.Context, user *entity.User, client *entity.Client) (*entity.Token, error)
	NewAccessToken(ctx context.Context, user *entity.User, client *entity.Client) (*entity.Token, error)
	NewRefreshToken(ctx context.Context, user *entity.User, client *entity.Client, ip, agent string) (*entity.Token, error)
	NewResetPassword(ctx context.Context, user *entity.User, client *entity.Client, ip, agent string) (*entity.Token, error)
}

type UseCase struct {
	repo   repository
	notify notify
	secure secure
}

func New(repo repository, notify notify, secure secure) *UseCase {
	return &UseCase{
		repo:   repo,
		notify: notify,
		secure: secure,
	}
}
