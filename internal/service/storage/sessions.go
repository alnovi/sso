package storage

import (
	"context"

	"github.com/alnovi/gomon/utils"
	"go.opentelemetry.io/otel/attribute"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/helper"
)

type Sessions struct {
	repo *repository.Repository
	tm   repository.Transaction
}

func NewSessions(repo *repository.Repository, tm repository.Transaction) *Sessions {
	return &Sessions{repo: repo, tm: tm}
}

func (s *Sessions) List(ctx context.Context) ([]*entity.SessionUser, error) {
	ctx, span := helper.SpanStart(ctx, "StorageSessions.List")
	defer span.End()

	session, err := s.repo.Sessions(ctx, repository.OrderDesc("updated_at"))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	userIds := utils.MapArray[string, *entity.Session](session, func(_ int, session *entity.Session) string {
		return session.UserId
	})

	users, err := s.repo.UserByIds(ctx, userIds)
	if err != nil {
		helper.SpanError(span, err)
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
	ctx, span := helper.SpanStart(ctx, "StorageSessions.GetById", helper.SpanAttr(
		attribute.String("session.id", id),
	))
	defer span.End()

	session, err := s.repo.SessionById(ctx, id)
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	user, err := s.repo.UserById(ctx, session.UserId)
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return &entity.SessionUser{Session: session, User: user}, nil
}

func (s *Sessions) DeleteById(ctx context.Context, id string) error {
	ctx, span := helper.SpanStart(ctx, "StorageSessions.DeleteById", helper.SpanAttr(
		attribute.String("session.id", id),
	))
	defer span.End()

	err := s.repo.SessionDeleteById(ctx, id)
	helper.SpanError(span, err)

	return err
}
