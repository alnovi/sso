package crontask

import (
	"context"

	"github.com/alnovi/sso/internal/adapter/repository"
)

type TaskDeleteSessionEmpty struct {
	repo *repository.Repository
}

func NewTaskDeleteSessionEmpty(repo *repository.Repository) *TaskDeleteSessionEmpty {
	return &TaskDeleteSessionEmpty{repo: repo}
}

func (t *TaskDeleteSessionEmpty) Handle() error {
	return t.repo.SessionDeleteWithoutTokens(context.Background())
}
