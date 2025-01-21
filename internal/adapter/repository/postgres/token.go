package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
)

const TokenTable = "tokens"

var tokenFields = []string{
	"id",
	"class",
	"hash",
	"session_id",
	"user_id",
	"client_id",
	"payload",
	"not_before",
	"expiration",
	"created_at",
	"updated_at",
}

func (r *Repository) TokenCreate(ctx context.Context, token *entity.Token) error {
	now := time.Now()

	if token.Id == "" {
		token.Id = uuid.NewString()
	}

	if token.NotBefore.IsZero() {
		token.NotBefore = now
	}

	if token.Expiration.IsZero() {
		token.Expiration = now
	}

	if token.CreatedAt.IsZero() {
		token.CreatedAt = now
	}

	if token.UpdatedAt.IsZero() {
		token.UpdatedAt = now
	}

	if token.Payload == nil {
		token.Payload = entity.Payload{}
	}

	query := r.qb.Insert(TokenTable).
		Columns(tokenFields...).
		Values(
			token.Id,
			token.Class,
			token.Hash,
			token.SessionId,
			token.UserId,
			token.ClientId,
			token.Payload,
			token.NotBefore,
			token.Expiration,
			token.CreatedAt,
			token.UpdatedAt,
		)

	q, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, q, args...)

	return err
}

func (r *Repository) TokenDelete(ctx context.Context, id string) error {
	query := r.qb.Delete(TokenTable).Where(sq.Eq{"id": id})

	q, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, q, args...)

	return err
}

func (r *Repository) TokenById(ctx context.Context, id string, fu bool) (*entity.Token, error) {
	token := new(entity.Token)

	err := uuid.Validate(id)
	if err != nil {
		return nil, repository.ErrNoResults
	}

	query := r.qb.Select(tokenFields...).
		From(TokenTable).
		Where(sq.Eq{"id": id})

	if fu {
		query = query.Suffix("FOR UPDATE")
	}

	q, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQueryRow(ctx, token, q, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNoResults
	}

	return token, err
}

func (r *Repository) TokenByClassHash(ctx context.Context, class, hash string, fu bool) (*entity.Token, error) {
	token := entity.NewToken()

	query := r.qb.Select(tokenFields...).
		From(TokenTable).
		Where(sq.Eq{"class": class, "hash": hash})

	if fu {
		query = query.Suffix("FOR UPDATE")
	}

	q, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQueryRow(ctx, token, q, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNoResults
	}

	return token, err
}
