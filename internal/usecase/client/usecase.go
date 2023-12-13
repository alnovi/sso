package client

import (
	"context"

	"github.com/alnovi/sso/internal/dto"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
	"github.com/alnovi/sso/internal/service"
)

type UseCase struct {
	client service.Client
}

func New(client service.Client) *UseCase {
	return &UseCase{client: client}
}

func (uc *UseCase) ClientForAuth(ctx context.Context, dto dto.ClientForAuth) (*entity.Client, error) {
	if dto.ClientId == "" {
		return uc.client.GetProfileClient(ctx)
	}

	client, err := uc.client.GetClient(ctx, dto.ClientId, &dto.RedirectURI, nil)
	if err != nil {
		return nil, exception.Wrap(exception.ClientNotFound, err)
	}

	return client, nil
}

func (uc *UseCase) ClientForToken(ctx context.Context, dto dto.ClientForToken) (*entity.Client, error) {
	client, err := uc.client.GetClient(ctx, dto.ClientId, nil, &dto.ClientSecret)
	if err != nil {
		return nil, exception.Wrap(exception.ClientNotFound, err)
	}

	return client, nil
}
