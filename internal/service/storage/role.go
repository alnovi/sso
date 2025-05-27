package storage

import (
	"context"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/pkg/utils"
)

type Roles struct {
	repo *repository.Repository
	tm   repository.Transaction
}

func NewRoles(repo *repository.Repository, tm repository.Transaction) *Roles {
	return &Roles{repo: repo, tm: tm}
}

func (s *Roles) ClientRoleByUserId(ctx context.Context, userId string) ([]*entity.ClientRole, error) {
	mapClientRole := make(map[string]*string)

	roles, err := s.repo.RoleByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		mapClientRole[role.ClientId] = &role.Role
	}

	clients, err := s.repo.Clients(ctx, repository.OrderAsc("name"), repository.NotDeleted())
	if err != nil {
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
	if userRole == nil || *userRole == "" {
		return s.repo.RoleDelete(ctx, clientId, userId)
	}

	return s.repo.RoleUpdate(ctx, &entity.Role{
		ClientId: clientId,
		UserId:   userId,
		Role:     *userRole,
	})
}
