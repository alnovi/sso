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

const SessionTable = "sessions"

var sessionFields = []string{"id", "user_id", "ip", "agent", "created_at", "updated_at"}

func (r *Repository) Sessions(ctx context.Context, opts ...OptSelect) ([]*entity.Session, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.Sessions")
	defer span.End()

	session := make([]*entity.Session, 0)

	builder := r.qb.Select(sessionFields...).From(SessionTable)
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQuery(ctx, &session, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return session, nil
}

func (r *Repository) SessionsCount(ctx context.Context, opts ...OptSelect) (int, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.SessionsCount")
	defer span.End()

	count := 0

	builder := r.qb.Select("COUNT (*)").From(SessionTable)
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return count, err
	}

	err = r.checkErr(r.db.QueryRow(ctx, query, args...).Scan(&count))
	if err != nil {
		helper.SpanError(span, err)
		return count, err
	}

	return count, nil
}

func (r *Repository) SessionById(ctx context.Context, id string) (*entity.Session, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.SessionById", helper.SpanAttr(
		attribute.String("session.id", id),
	))
	defer span.End()

	session := new(entity.Session)

	if err := r.checkUUID(id); err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	builder := r.qb.Select(sessionFields...).
		From(SessionTable).
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQueryRow(ctx, session, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return session, nil
}

func (r *Repository) SessionByUserId(ctx context.Context, userId string, opts ...OptSelect) (*entity.Session, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.SessionByUserId", helper.SpanAttr(
		attribute.String("user.id", userId),
	))
	defer span.End()

	session := new(entity.Session)

	if err := r.checkUUID(userId); err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	builder := r.qb.Select(sessionFields...).
		From(SessionTable).
		Where(sq.Eq{"user_id": userId}).
		Limit(1)
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQueryRow(ctx, session, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return session, nil
}

func (r *Repository) SessionsByUserId(ctx context.Context, userId string, opts ...OptSelect) ([]*entity.Session, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.SessionsByUserId", helper.SpanAttr(
		attribute.String("user.id", userId),
	))
	defer span.End()

	sessions := make([]*entity.Session, 0)

	if err := r.checkUUID(userId); err != nil {
		return nil, err
	}

	builder := r.qb.Select(sessionFields...).
		From(SessionTable).
		Where(sq.Eq{"user_id": userId})

	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQuery(ctx, &sessions, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return sessions, nil
}

func (r *Repository) SessionCreate(ctx context.Context, session *entity.Session) error {
	ctx, span := helper.SpanStart(ctx, "Repository.SessionCreate")
	defer span.End()

	now := time.Now()

	if session.Id == "" {
		session.Id = uuid.NewString()
	}

	if session.CreatedAt.IsZero() {
		session.CreatedAt = now
	}

	if session.UpdatedAt.IsZero() {
		session.UpdatedAt = now
	}

	span.SetAttributes(attribute.String("session.id", session.Id))

	builder := r.qb.Insert(SessionTable).
		Columns(sessionFields...).
		Values(
			session.Id,
			session.UserId,
			session.Ip,
			session.Agent,
			session.CreatedAt,
			session.UpdatedAt,
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

func (r *Repository) SessionUpdateDateById(ctx context.Context, id string) error {
	ctx, span := helper.SpanStart(ctx, "Repository.SessionUpdateDateById", helper.SpanAttr(
		attribute.String("session.id", id),
	))
	defer span.End()

	if err := r.checkUUID(id); err != nil {
		helper.SpanError(span, err)
		return err
	}

	builder := r.qb.Update(SessionTable).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id})

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

func (r *Repository) SessionDeleteById(ctx context.Context, id string) error {
	ctx, span := helper.SpanStart(ctx, "Repository.SessionDeleteById", helper.SpanAttr(
		attribute.String("session.id", id),
	))
	defer span.End()

	if err := r.checkUUID(id); err != nil {
		helper.SpanError(span, err)
		return err
	}

	builder := r.qb.Delete(SessionTable).Where(sq.Eq{"id": id})

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

func (r *Repository) SessionDeleteWithoutTokens(ctx context.Context) error {
	ctx, span := helper.SpanStart(ctx, "Repository.SessionDeleteWithoutTokens")
	defer span.End()

	builder := r.qb.Delete(SessionTable).Where(sq.Expr("NOT EXISTS(SELECT * FROM tokens WHERE session_id = sessions.id)"))

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
