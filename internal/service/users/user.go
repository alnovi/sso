package users

import (
	"context"
	"errors"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type User struct {
	repo repository.Repository
}

func New(repo repository.Repository) *User {
	return &User{repo: repo}
}

func (s *User) UserById(ctx context.Context, id string) (*entity.User, error) {
	user, err := s.repo.UserById(ctx, id)
	if err != nil && !errors.Is(err, repository.ErrNoResults) {
		return nil, err
	}
	return user, err
}
