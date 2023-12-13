package user

import (
	"context"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CanUseClient(ctx context.Context, user entity.User, client entity.Client) error {
	return s.repo.CanUseClient(ctx, client, user)
}

func (s *Service) GetUserById(ctx context.Context, userId string) (*entity.User, error) {
	return s.repo.GetUserById(ctx, userId)
}
