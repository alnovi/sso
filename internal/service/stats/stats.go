package stats

import (
	"context"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/helper"
)

type Stats struct {
	repo *repository.Repository
}

func NewStats(repo *repository.Repository) *Stats {
	return &Stats{repo: repo}
}

func (s *Stats) UserCount(ctx context.Context) (int, error) {
	ctx, span := helper.SpanStart(ctx, "Stats.UserCount")
	defer span.End()

	count, err := s.repo.UsersCount(ctx, repository.NotDeleted())
	helper.SpanError(span, err)

	return count, err
}

func (s *Stats) ClientCount(ctx context.Context) (int, error) {
	ctx, span := helper.SpanStart(ctx, "Stats.ClientCount")
	defer span.End()

	count, err := s.repo.ClientsCount(ctx, repository.NotDeleted())
	helper.SpanError(span, err)

	return count, err
}

func (s *Stats) SessionCount(ctx context.Context) (int, error) {
	ctx, span := helper.SpanStart(ctx, "Stats.SessionCount")
	defer span.End()

	count, err := s.repo.SessionsCount(ctx)
	helper.SpanError(span, err)

	return count, err
}
