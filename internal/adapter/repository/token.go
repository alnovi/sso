package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/helper"
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
	ctx, span := helper.SpanStart(ctx, "Repository.TokenByHash")
	defer span.End()

	token := new(entity.Token)

	builder := r.qb.Select(tokenFields...).
		From(TokenTable).
		Where(sq.Eq{"hash": hash})

	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQueryRow(ctx, token, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return token, nil
}

func (r *Repository) TokenCreate(ctx context.Context, token *entity.Token) error {
	ctx, span := helper.SpanStart(ctx, "Repository.TokenCreate")
	defer span.End()

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

	span.SetAttributes(attribute.String("token.id", token.Id))

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
		helper.SpanError(span, err)
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err = r.checkErr(err); err != nil {
		helper.SpanError(span, err)
		return err
	}

	return nil
}

func (r *Repository) TokenDeleteById(ctx context.Context, id string) error {
	ctx, span := helper.SpanStart(ctx, "Repository.TokenDeleteById", helper.SpanAttr(
		attribute.String("token.id", id),
	))
	defer span.End()

	builder := r.qb.Delete(TokenTable).Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err = r.checkErr(err); err != nil {
		helper.SpanError(span, err)
		return err
	}

	return nil
}

func (r *Repository) TokenDeleteBySessionId(ctx context.Context, sessionId string) error {
	ctx, span := helper.SpanStart(ctx, "Repository.TokenDeleteBySessionId", helper.SpanAttr(
		attribute.String("session.id", sessionId),
	))
	defer span.End()

	if err := r.checkUUID(sessionId); err != nil {
		helper.SpanError(span, err)
		return err
	}

	builder := r.qb.Delete(TokenTable).Where(sq.Eq{"session_id": sessionId})

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err = r.checkErr(err); err != nil {
		helper.SpanError(span, err)
		return err
	}

	return nil
}

func (r *Repository) TokenDeleteExpired(ctx context.Context) error {
	ctx, span := helper.SpanStart(ctx, "Repository.TokenDeleteExpired")
	defer span.End()

	builder := r.qb.Delete(TokenTable).
		Where(sq.NotEq{"expiration": nil}).
		Where(sq.Lt{"expiration": time.Now()})

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	if err = r.checkErr(err); err != nil {
		helper.SpanError(span, err)
		return err
	}

	return nil
}
