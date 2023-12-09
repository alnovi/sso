package client

import (
	"context"
	"errors"
	"strings"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetProfileClient(ctx context.Context) (*entity.Client, error) {
	return s.repo.GetClientByClass(ctx, entity.ClientClassProfile)
}

func (s *Service) GetClient(ctx context.Context, id string, callback, secret *string) (*entity.Client, error) {
	client, err := s.repo.GetClientByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if callback != nil && strings.Trim(*callback, "/") != strings.Trim(client.Callback, "/") {
		return nil, errors.New("client is not match callback")
	}

	if secret != nil && client.Secret != *secret {
		return nil, errors.New("client is not match secret")
	}

	return client, nil
}
