package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/alnovi/sso/pkg/database/postgres"
)

var (
	ErrNoResult        = errors.New("no results")
	ErrClientIdExists  = errors.New("client id exists")
	ErrUserEmailExists = errors.New("user email exists")
)

type Transaction interface {
	ReadCommitted(ctx context.Context, fn func(ctx context.Context) error) error
}

type Repository struct {
	db *postgres.Client
	qb sq.StatementBuilderType
}

func NewRepository(db *postgres.Client) *Repository {
	return &Repository{
		db: db,
		qb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *Repository) checkUUID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return fmt.Errorf("%w: invalid uuid", ErrNoResult)
	}
	return nil
}

func (r *Repository) checkErr(err error) error {
	var pgErr *pgconn.PgError

	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoResult
	}

	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" && pgErr.ConstraintName == "users_email_unique" {
			return ErrUserEmailExists
		}

		if pgErr.Code == "23505" && pgErr.ConstraintName == "clients_pkey" {
			return ErrClientIdExists
		}
	}

	return err
}

func (r *Repository) applyOptSelect(query sq.SelectBuilder, opts []OptSelect) sq.SelectBuilder {
	for _, opt := range opts {
		query = opt(query)
	}
	return query
}
