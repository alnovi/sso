package storage

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel/attribute"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/helper"
	"github.com/alnovi/sso/pkg/rand"
)

const secretLength = 50

var (
	ErrClientIdExists = errors.New("client id exists")
)

type Clients struct {
	repo *repository.Repository
	tm   repository.Transaction
}

func NewClients(repo *repository.Repository, tm repository.Transaction) *Clients {
	return &Clients{repo: repo, tm: tm}
}

func (s *Clients) All(ctx context.Context) ([]*entity.Client, error) {
	ctx, span := helper.SpanStart(ctx, "StorageClients.All")
	defer span.End()

	clients, err := s.repo.Clients(ctx, repository.OrderAsc("name"))
	helper.SpanError(span, err)

	return clients, err
}

func (s *Clients) GetById(ctx context.Context, id string) (*entity.Client, error) {
	ctx, span := helper.SpanStart(ctx, "StorageClients.GetById", helper.SpanAttr(
		attribute.String("client.id", id),
	))
	defer span.End()

	client, err := s.repo.ClientById(ctx, id)
	helper.SpanError(span, err)

	return client, err
}

func (s *Clients) Create(ctx context.Context, inp InputClientCreate) (*entity.Client, error) {
	ctx, span := helper.SpanStart(ctx, "StorageClients.Create", helper.SpanAttr(
		attribute.String("client.id", inp.Id),
		attribute.String("client.name", inp.Name),
	))
	defer span.End()

	if inp.Secret == nil || *inp.Secret == "" {
		secret := rand.Base62(secretLength)
		inp.Secret = &secret
	}

	client := &entity.Client{
		Id:       inp.Id,
		Name:     inp.Name,
		Icon:     inp.Icon,
		Secret:   *inp.Secret,
		Callback: inp.Callback,
		IsSystem: false,
	}

	err := s.checkErr(s.repo.ClientCreate(ctx, client))
	helper.SpanError(span, err)

	return client, err
}

func (s *Clients) Update(ctx context.Context, inp InputClientUpdate) (*entity.Client, error) {
	ctx, span := helper.SpanStart(ctx, "StorageClients.Update", helper.SpanAttr(
		attribute.String("client.id", inp.Id),
	))
	defer span.End()

	client, err := s.GetById(ctx, inp.Id)
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	client.Name = inp.Name
	client.Icon = inp.Icon
	client.Callback = inp.Callback
	client.Secret = inp.Secret

	err = s.checkErr(s.repo.ClientUpdate(ctx, client))
	helper.SpanError(span, err)

	return client, err
}

func (s *Clients) Delete(ctx context.Context, id string) (*entity.Client, error) {
	ctx, span := helper.SpanStart(ctx, "StorageClients.Delete", helper.SpanAttr(
		attribute.String("client.id", id),
	))
	defer span.End()

	client, err := s.repo.ClientById(ctx, id, repository.NotSystem())
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	if client.DeletedAt == nil {
		err = s.repo.ClientDelete(ctx, client)
	} else {
		err = s.repo.ClientDeleteForce(ctx, client)
	}

	helper.SpanError(span, err)

	return client, err
}

func (s *Clients) Restore(ctx context.Context, id string) (*entity.Client, error) {
	ctx, span := helper.SpanStart(ctx, "StorageClients.Restore", helper.SpanAttr(
		attribute.String("client.id", id),
	))
	defer span.End()

	client, err := s.repo.ClientById(ctx, id, repository.IsNotNull("deleted_at"))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	client.DeletedAt = nil

	err = s.repo.ClientUpdate(ctx, client)
	helper.SpanError(span, err)

	return client, err
}

func (s *Clients) checkErr(err error) error {
	if errors.Is(err, repository.ErrClientIdExists) {
		return ErrClientIdExists
	}
	return err
}
