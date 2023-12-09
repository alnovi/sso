package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
	"github.com/google/uuid"
)

var tokenFields = []string{
	"id",
	"class",
	"hash",
	"user_id",
	"client_id",
	"meta",
	"not_before",
	"expiration",
	"created_at",
	"updated_at",
}

func (r *Repository) CreateToken(ctx context.Context, token *entity.Token) error {
	now := time.Now()

	token.Id = uuid.NewString()
	token.CreatedAt = now
	token.UpdatedAt = now

	if _, err := uuid.Parse(*token.UserId); err != nil {
		token.UserId = nil
	}

	if _, err := uuid.Parse(*token.ClientId); err != nil {
		token.ClientId = nil
	}

	_, err := r.qb.Insert(tableTokens).
		Columns(tokenFields...).
		Values(
			token.Id,
			token.Class,
			token.Hash,
			token.UserId,
			token.ClientId,
			token.Meta,
			token.NotBefore,
			token.Expiration,
			token.CreatedAt,
			token.UpdatedAt,
		).
		RunWith(r.connect(ctx)).
		Exec()

	return err
}

func (r *Repository) UpdateToken(ctx context.Context, token *entity.Token) error {
	token.UpdatedAt = time.Now()

	if _, err := uuid.Parse(*token.UserId); err != nil {
		token.UserId = nil
	}

	if _, err := uuid.Parse(*token.ClientId); err != nil {
		token.ClientId = nil
	}

	_, err := r.qb.Update(tableTokens).
		Set("hash", token.Hash).
		Set("user_id", token.UserId).
		Set("client_id", token.ClientId).
		Set("meta", token.Meta).
		Set("meta", token.NotBefore).
		Set("meta", token.Expiration).
		Set("meta", token.UpdatedAt).
		Where(squirrel.Eq{"id": token.Id}).
		RunWith(r.connect(ctx)).
		Exec()

	return err
}

func (r *Repository) GetTokenByHash(ctx context.Context, hash string) (*entity.Token, error) {
	result := &entity.Token{}

	err := r.qb.Select(tokenFields...).
		From(tableTokens).
		Where(squirrel.Eq{"hash": hash}).
		RunWith(r.connect(ctx)).
		QueryRow().
		Scan(
			&result.Id,
			&result.Class,
			&result.Hash,
			&result.UserId,
			&result.ClientId,
			&result.Meta,
			&result.NotBefore,
			&result.Expiration,
			&result.CreatedAt,
			&result.UpdatedAt,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, exception.TokenNotFound
	}

	return result, err
}
