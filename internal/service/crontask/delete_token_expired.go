package crontask

import (
	"context"

	"github.com/alnovi/sso/internal/adapter/repository"
)

type TaskDeleteTokenExpired struct {
	repo *repository.Repository
}

func NewTaskDeleteTokenExpired(repo *repository.Repository) *TaskDeleteTokenExpired {
	return &TaskDeleteTokenExpired{repo: repo}
}

func (t *TaskDeleteTokenExpired) Handle() error {
	return t.repo.TokenDeleteExpired(context.Background())
}
