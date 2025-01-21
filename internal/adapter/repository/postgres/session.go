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

const SessionTable = "sessions"

var sessionFields = []string{
	"id",
	"user_id",
	"ip",
	"agent",
	"created_at",
	"updated_at",
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

	query := r.qb.Insert(SessionTable).
		Columns(sessionFields...).
		Values(
			session.Id,
			session.UserId,
			session.Ip,
			session.Agent,
			session.CreatedAt,
			session.UpdatedAt,
		)

	q, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, q, args...)

	return err
}

func (r *Repository) SessionDelete(ctx context.Context, id string) error {
	err := uuid.Validate(id)
	if err != nil {
		return repository.ErrNoResults
	}

	query := r.qb.Delete(SessionTable).Where(sq.Eq{"id": id})

	q, args, err := query.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, q, args...)

	return err
}

func (r *Repository) SessionById(ctx context.Context, id string) (*entity.Session, error) {
	session := new(entity.Session)

	err := uuid.Validate(id)
	if err != nil {
		return nil, repository.ErrNoResults
	}

	query := r.qb.Select(sessionFields...).
		From(SessionTable).
		Where(sq.Eq{"id": id})

	q, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQueryRow(ctx, session, q, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNoResults
	}

	return session, err
}
