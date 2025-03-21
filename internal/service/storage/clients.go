package storage

import (
	"context"
	"errors"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
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
	return s.repo.Clients(ctx, repository.OrderAsc("name"))
}

func (s *Clients) GetById(ctx context.Context, id string) (*entity.Client, error) {
	return s.repo.ClientById(ctx, id)
}

func (s *Clients) Create(ctx context.Context, inp InputClientCreate) (*entity.Client, error) {
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

	err := s.repo.ClientCreate(ctx, client)

	return client, s.checkErr(err)
}

func (s *Clients) Update(ctx context.Context, inp InputClientUpdate) (*entity.Client, error) {
	client, err := s.GetById(ctx, inp.Id)
	if err != nil {
		return nil, err
	}

	client.Name = inp.Name
	client.Icon = inp.Icon
	client.Callback = inp.Callback
	client.Secret = inp.Secret

	err = s.repo.ClientUpdate(ctx, client)

	return client, s.checkErr(err)
}

func (s *Clients) Delete(ctx context.Context, id string) (*entity.Client, error) {
	client, err := s.repo.ClientById(ctx, id, repository.NotSystem())
	if err != nil {
		return nil, err
	}

	if client.DeletedAt == nil {
		err = s.repo.ClientDelete(ctx, client)
	} else {
		err = s.repo.ClientDeleteForce(ctx, client)
	}

	return client, err
}

func (s *Clients) Restore(ctx context.Context, id string) (*entity.Client, error) {
	client, err := s.repo.ClientById(ctx, id, repository.IsNotNull("deleted_at"))
	if err != nil {
		return nil, err
	}

	client.DeletedAt = nil

	return client, s.repo.ClientUpdate(ctx, client)
}

func (s *Clients) checkErr(err error) error {
	if errors.Is(err, repository.ErrClientIdExists) {
		return ErrClientIdExists
	}
	return err
}
