package storage

import (
	"context"
	"errors"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/pkg/utils"
)

var (
	ErrUserEmailExists = errors.New("user email exists")
)

type Users struct {
	repo *repository.Repository
	tm   repository.Transaction
}

func NewUsers(repo *repository.Repository, tm repository.Transaction) *Users {
	return &Users{repo: repo, tm: tm}
}

func (s *Users) All(ctx context.Context) ([]*entity.User, error) {
	return s.repo.Users(ctx, repository.OrderAsc("name"))
}

func (s *Users) GetById(ctx context.Context, id string) (*entity.User, error) {
	return s.repo.UserById(ctx, id)
}

func (s *Users) Create(ctx context.Context, inp InputUserCreate) (*entity.User, error) {
	password, err := utils.HashPassword(inp.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Name:     inp.Name,
		Email:    inp.Email,
		Password: password,
	}

	err = s.repo.UserCreate(ctx, user)

	return user, s.checkErr(err)
}

func (s *Users) Update(ctx context.Context, inp InputUserUpdate) (*entity.User, error) {
	user, err := s.GetById(ctx, inp.Id)
	if err != nil {
		return nil, err
	}

	user.Name = inp.Name
	user.Email = inp.Email

	if inp.Password != nil {
		user.Password, err = utils.HashPassword(*inp.Password)
		if err != nil {
			return nil, err
		}
	}

	err = s.repo.UserUpdate(ctx, user)

	return user, s.checkErr(err)
}

func (s *Users) Delete(ctx context.Context, id string) (*entity.User, error) {
	user, err := s.repo.UserById(ctx, id)
	if err != nil {
		return nil, err
	}

	if user.DeletedAt == nil {
		err = s.repo.UserDelete(ctx, user)
	} else {
		err = s.repo.UserDeleteForce(ctx, user)
	}

	return user, err
}

func (s *Users) Restore(ctx context.Context, id string) (*entity.User, error) {
	user, err := s.repo.UserById(ctx, id, repository.IsNotNull("deleted_at"))
	if err != nil {
		return nil, err
	}

	user.DeletedAt = nil

	return user, s.repo.UserUpdate(ctx, user)
}

func (s *Users) checkErr(err error) error {
	if errors.Is(err, repository.ErrUserEmailExists) {
		return ErrUserEmailExists
	}
	return err
}
