package stats

import (
	"context"

	"github.com/alnovi/sso/internal/adapter/repository"
)

type Stats struct {
	repo *repository.Repository
}

func NewStats(repo *repository.Repository) *Stats {
	return &Stats{repo: repo}
}

func (s *Stats) UserCount(ctx context.Context) (int, error) {
	return s.repo.UsersCount(ctx, repository.NotDeleted())
}

func (s *Stats) ClientCount(ctx context.Context) (int, error) {
	return s.repo.ClientsCount(ctx, repository.NotDeleted())
}

func (s *Stats) SessionCount(ctx context.Context) (int, error) {
	return s.repo.SessionsCount(ctx)
}
