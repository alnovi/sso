package repository

import (
	"context"
	"errors"

	"github.com/alnovi/sso/internal/entity"
)

var (
	ErrNoResults = errors.New("no results")
)

type Transaction interface {
	ReadCommitted(ctx context.Context, fn func(ctx context.Context) error) error
}

type Repository interface {
	ClientById(ctx context.Context, id string) (*entity.Client, error)

	UserByEmail(ctx context.Context, email string) (*entity.User, error)

	SessionCreate(ctx context.Context, session *entity.Session) error
	SessionDelete(ctx context.Context, id string) error
	SessionById(ctx context.Context, id string) (*entity.Session, error)

	TokenCreate(ctx context.Context, token *entity.Token) error
	TokenDelete(ctx context.Context, id string) error
	TokenById(ctx context.Context, id string, fu bool) (*entity.Token, error)
	TokenByClassHash(ctx context.Context, class, hash string, fu bool) (*entity.Token, error)
}
