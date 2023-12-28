package usecase

import (
	"context"
	"net/url"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/entity"
)

type Auth interface {
	AuthById(ctx context.Context, dto dto.AuthById) (*entity.User, *url.URL, error)
	AuthByCredentials(ctx context.Context, dto dto.AuthByCredentials) (*entity.User, *url.URL, error)
}

type Client interface {
	ClientForAuth(ctx context.Context, dto dto.ClientForAuth) (*entity.Client, error)
	ClientForToken(ctx context.Context, dto dto.ClientForToken) (*entity.Client, error)
}

type Token interface {
	AccessAndRefreshToken(ctx context.Context, dto dto.AccessToken) (*entity.TokenWithUser, *entity.Token, error)
}

type Profile interface {
	Profile(ctx context.Context, userId string) (*entity.User, []*entity.Token, []*entity.Client, error)
}
