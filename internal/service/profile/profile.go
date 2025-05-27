package profile

import (
	"context"
	"errors"
	"fmt"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/pkg/utils"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

type UserProfile struct {
	repo *repository.Repository
	tm   repository.Transaction
}

func NewUserProfile(repo *repository.Repository, tm repository.Transaction) *UserProfile {
	return &UserProfile{repo: repo, tm: tm}
}

func (s *UserProfile) SessionByIdAndAgent(ctx context.Context, id, agent string) (*entity.Session, error) {
	session, err := s.repo.SessionById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrSessionNotFound, err)
	}

	if session.Agent != agent {
		return nil, fmt.Errorf("%w: agent not attempted", ErrSessionNotFound)
	}

	return s.repo.SessionById(ctx, id)
}

func (s *UserProfile) Info(ctx context.Context, userId string) (*entity.User, error) {
	user, err := s.repo.UserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUserNotFound, err)
	}
	return user, nil
}

func (s *UserProfile) UpdateInfo(ctx context.Context, userId, name, email string) (*entity.User, error) {
	user, err := s.repo.UserById(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrUserNotFound, err)
	}

	user.Name = name
	user.Email = email

	err = s.repo.UserUpdate(ctx, user)

	return user, err
}

func (s *UserProfile) Clients(ctx context.Context, userId string) ([]*entity.ClientRole, error) {
	mapClientRole := make(map[string]*string)
	clientIds := make([]string, 0)

	roles, err := s.repo.RoleByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	for _, role := range roles {
		mapClientRole[role.ClientId] = &role.Role
		clientIds = append(clientIds, role.ClientId)
	}

	clients, err := s.repo.ClientByIds(ctx, clientIds, repository.OrderAsc("name"))
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

func (s *UserProfile) Sessions(ctx context.Context, userId string) ([]*entity.Session, error) {
	return s.repo.SessionsByUserId(ctx, userId, repository.OrderDesc("created_at"))
}

func (s *UserProfile) SessionDelete(ctx context.Context, userId, sessionId string) error {
	session, err := s.repo.SessionById(ctx, sessionId)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrSessionNotFound, err)
	}

	if session.UserId != userId {
		return fmt.Errorf("%w: session not attempted", ErrSessionNotFound)
	}

	return s.repo.SessionDeleteById(ctx, sessionId)
}

func (s *UserProfile) UpdatePassword(ctx context.Context, userId, oldPassword, newPassword string) error {
	var err error

	user, err := s.repo.UserById(ctx, userId)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrUserNotFound, err)
	}

	if !utils.CompareHashPassword(oldPassword, user.Password) {
		return ErrInvalidPassword
	}

	if user.Password, err = utils.HashPassword(newPassword); err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidPassword, err)
	}

	return s.repo.UserUpdate(ctx, user)
}

func (s *UserProfile) Logout(ctx context.Context, sessionId string) error {
	return s.repo.SessionDeleteById(ctx, sessionId)
}
