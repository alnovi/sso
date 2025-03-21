package storage

import (
	"context"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/pkg/utils"
)

type Sessions struct {
	repo *repository.Repository
	tm   repository.Transaction
}

func NewSessions(repo *repository.Repository, tm repository.Transaction) *Sessions {
	return &Sessions{repo: repo, tm: tm}
}

func (s *Sessions) List(ctx context.Context) ([]*entity.SessionUser, error) {
	session, err := s.repo.Sessions(ctx, repository.OrderDesc("updated_at"))
	if err != nil {
		return nil, err
	}

	userIds := utils.MapArray[string, *entity.Session](session, func(_ int, session *entity.Session) string {
		return session.UserId
	})

	users, err := s.repo.UserByIds(ctx, userIds)
	if err != nil {
		return nil, err
	}

	return utils.MapArray[*entity.SessionUser, *entity.Session](session, func(_ int, session *entity.Session) *entity.SessionUser {
		for _, user := range users {
			if user.Id == session.UserId {
				return &entity.SessionUser{
					Session: session,
					User:    user,
				}
			}
		}
		return &entity.SessionUser{Session: session}
	}), nil
}

func (s *Sessions) GetById(ctx context.Context, id string) (*entity.SessionUser, error) {
	session, err := s.repo.SessionById(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.UserById(ctx, session.UserId)
	if err != nil {
		return nil, err
	}

	return &entity.SessionUser{Session: session, User: user}, nil
}

func (s *Sessions) DeleteById(ctx context.Context, id string) error {
	return s.repo.SessionDeleteById(ctx, id)
}
