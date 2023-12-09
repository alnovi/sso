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
}
