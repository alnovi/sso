package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/alnovi/sso/internal/entity"
)

const SessionTable = "sessions"

var sessionFields = []string{"id", "user_id", "ip", "agent", "created_at", "updated_at"}

func (r *Repository) Sessions(ctx context.Context, opts ...OptSelect) ([]*entity.Session, error) {
	session := make([]*entity.Session, 0)

	builder := r.qb.Select(sessionFields...).From(SessionTable)
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQuery(ctx, &session, query, args...)

	return session, r.checkErr(err)
}

func (r *Repository) SessionsCount(ctx context.Context, opts ...OptSelect) (int, error) {
	count := 0

	builder := r.qb.Select("COUNT (*)").From(SessionTable)
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return count, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&count)

	return count, r.checkErr(err)
}

func (r *Repository) SessionById(ctx context.Context, id string) (*entity.Session, error) {
	session := new(entity.Session)

	if err := r.checkUUID(id); err != nil {
		return nil, err
	}

	builder := r.qb.Select(sessionFields...).
		From(SessionTable).
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQueryRow(ctx, session, query, args...)

	return session, r.checkErr(err)
}

func (r *Repository) SessionByUserId(ctx context.Context, userId string, opts ...OptSelect) (*entity.Session, error) {
	session := new(entity.Session)

	if err := r.checkUUID(userId); err != nil {
		return nil, err
	}

	builder := r.qb.Select(sessionFields...).
		From(SessionTable).
		Where(sq.Eq{"user_id": userId}).
		Limit(1)
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQueryRow(ctx, session, query, args...)

	return session, r.checkErr(err)
}

func (r *Repository) SessionsByUserId(ctx context.Context, userId string, opts ...OptSelect) ([]*entity.Session, error) {
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
		return nil, err
	}

	err = r.db.ScanQuery(ctx, &sessions, query, args...)

	return sessions, r.checkErr(err)
}

func (r *Repository) SessionCreate(ctx context.Context, session *entity.Session) error {
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
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) SessionUpdateDateById(ctx context.Context, id string) error {
	if err := r.checkUUID(id); err != nil {
		return err
	}

	builder := r.qb.Update(SessionTable).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) SessionDeleteById(ctx context.Context, id string) error {
	if err := r.checkUUID(id); err != nil {
		return err
	}

	builder := r.qb.Delete(SessionTable).Where(sq.Eq{"id": id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) SessionDeleteByUserId(ctx context.Context, userId string) error {
	if err := r.checkUUID(userId); err != nil {
		return err
	}

	builder := r.qb.Delete(SessionTable).Where(sq.Eq{"user_id": userId})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) SessionDeleteWithoutTokens(ctx context.Context) error {
	builder := r.qb.Delete(SessionTable).Where(sq.Expr("NOT EXISTS(SELECT * FROM tokens WHERE session_id = sessions.id)"))

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}
