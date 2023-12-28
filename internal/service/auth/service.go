package auth

import (
	"context"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
)

type Service struct {
	repo repository.Repository
}

func New(repo repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AuthByCredentials(ctx context.Context, login, password string) (*entity.User, error) {
	user, err := s.repo.GetUserByLoginOrEmail(ctx, login)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, exception.ErrPasswordIncorrect
	}

	return user, nil
}
