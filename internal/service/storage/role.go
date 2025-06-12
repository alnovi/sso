package storage

import (
	"context"

	"github.com/alnovi/gomon/utils"
	"go.opentelemetry.io/otel/attribute"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/helper"
)

type Roles struct {
	repo *repository.Repository
	tm   repository.Transaction
}

func NewRoles(repo *repository.Repository, tm repository.Transaction) *Roles {
	return &Roles{repo: repo, tm: tm}
}

func (s *Roles) ClientRoleByUserId(ctx context.Context, userId string) ([]*entity.ClientRole, error) {
	ctx, span := helper.SpanStart(ctx, "StorageRoles.ClientRoleByUserId", helper.SpanAttr(
		attribute.String("user.id", userId),
	))
	defer span.End()

	mapClientRole := make(map[string]*string)

	roles, err := s.repo.RoleByUserId(ctx, userId)
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	for _, role := range roles {
		mapClientRole[role.ClientId] = &role.Role
	}

	clients, err := s.repo.Clients(ctx, repository.OrderAsc("name"), repository.NotDeleted())
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return utils.MapArray[*entity.ClientRole, *entity.Client](clients, func(_ int, client *entity.Client) *entity.ClientRole {
		return &entity.ClientRole{
			Client: client,
			Role:   mapClientRole[client.Id],
		}
	}), nil
}

func (s *Roles) Update(ctx context.Context, clientId, userId string, userRole *string) error {
	ctx, span := helper.SpanStart(ctx, "StorageRoles.Update", helper.SpanAttr(
		attribute.String("client.id", clientId),
		attribute.String("user.id", userId),
	))
	defer span.End()

	if userRole == nil || *userRole == "" {
		err := s.repo.RoleDelete(ctx, clientId, userId)
		helper.SpanError(span, err)
		return err
	}

	err := s.repo.RoleUpdate(ctx, &entity.Role{
		ClientId: clientId,
		UserId:   userId,
		Role:     *userRole,
	})
	helper.SpanError(span, err)
	return err
}
