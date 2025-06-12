package storage

import (
	"context"
	"errors"

	"github.com/alnovi/gomon/utils"
	"go.opentelemetry.io/otel/attribute"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/helper"
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
	ctx, span := helper.SpanStart(ctx, "StorageUsers.All")
	defer span.End()

	users, err := s.repo.Users(ctx, repository.OrderAsc("name"))
	helper.SpanError(span, err)

	return users, err
}

func (s *Users) GetById(ctx context.Context, id string) (*entity.User, error) {
	ctx, span := helper.SpanStart(ctx, "StorageUsers.GetById", helper.SpanAttr(
		attribute.String("user.id", id),
	))
	defer span.End()

	user, err := s.repo.UserById(ctx, id)
	helper.SpanError(span, err)

	return user, err
}

func (s *Users) Create(ctx context.Context, inp InputUserCreate) (*entity.User, error) {
	ctx, span := helper.SpanStart(ctx, "StorageUsers.Create", helper.SpanAttr(
		attribute.String("user.email", inp.Email),
		attribute.String("user.name", inp.Name),
	))
	defer span.End()

	password, err := utils.HashPassword(inp.Password)
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	user := &entity.User{
		Name:     inp.Name,
		Email:    inp.Email,
		Password: password,
	}

	err = s.checkErr(s.repo.UserCreate(ctx, user))
	helper.SpanError(span, err)

	return user, err
}

func (s *Users) Update(ctx context.Context, inp InputUserUpdate) (*entity.User, error) {
	ctx, span := helper.SpanStart(ctx, "StorageUsers.Update", helper.SpanAttr(
		attribute.String("user.id", inp.Id),
		attribute.String("user.email", inp.Email),
		attribute.String("user.name", inp.Name),
	))
	defer span.End()

	user, err := s.GetById(ctx, inp.Id)
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	user.Name = inp.Name
	user.Email = inp.Email

	if inp.Password != nil {
		user.Password, err = utils.HashPassword(*inp.Password)
		if err != nil {
			helper.SpanError(span, err)
			return nil, err
		}
	}

	err = s.checkErr(s.repo.UserUpdate(ctx, user))
	helper.SpanError(span, err)

	return user, err
}

func (s *Users) Delete(ctx context.Context, id string) (*entity.User, error) {
	ctx, span := helper.SpanStart(ctx, "StorageUsers.Delete", helper.SpanAttr(
		attribute.String("user.id", id),
	))
	defer span.End()

	user, err := s.repo.UserById(ctx, id)
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	if user.DeletedAt == nil {
		err = s.repo.UserDelete(ctx, user)
	} else {
		err = s.repo.UserDeleteForce(ctx, user)
	}

	helper.SpanError(span, err)

	return user, err
}

func (s *Users) Restore(ctx context.Context, id string) (*entity.User, error) {
	ctx, span := helper.SpanStart(ctx, "StorageUsers.Restore", helper.SpanAttr(
		attribute.String("user.id", id),
	))
	defer span.End()

	user, err := s.repo.UserById(ctx, id, repository.IsNotNull("deleted_at"))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	user.DeletedAt = nil

	err = s.repo.UserUpdate(ctx, user)
	helper.SpanError(span, err)

	return user, err
}

func (s *Users) checkErr(err error) error {
	if errors.Is(err, repository.ErrUserEmailExists) {
		return ErrUserEmailExists
	}
	return err
}
