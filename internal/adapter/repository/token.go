package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

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

func (r *Repository) TokenByHash(ctx context.Context, hash string, opts ...OptSelect) (*entity.Token, error) {
	token := new(entity.Token)

	builder := r.qb.Select(tokenFields...).
		From(TokenTable).
		Where(sq.Eq{"hash": hash})

	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQueryRow(ctx, token, query, args...)

	return token, r.checkErr(err)
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

	builder := r.qb.Insert(TokenTable).
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

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) TokenDeleteById(ctx context.Context, id string) error {
	builder := r.qb.Delete(TokenTable).Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) TokenDeleteBySessionId(ctx context.Context, sessionId string) error {
	if err := r.checkUUID(sessionId); err != nil {
		return err
	}

	builder := r.qb.Delete(TokenTable).Where(sq.Eq{"session_id": sessionId})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) TokenDeleteExpired(ctx context.Context) error {
	builder := r.qb.Delete(TokenTable).
		Where(sq.NotEq{"expiration": nil}).
		Where(sq.Lt{"expiration": time.Now()})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}
