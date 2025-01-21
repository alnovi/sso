package postgres

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
)

const UserTable = "users"

var userFields = []string{"id", "name", "email", "password", "created_at", "updated_at"}

func (r *Repository) UserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := new(entity.User)

	query := r.qb.Select(userFields...).
		From(UserTable).
		Where(sq.Eq{"email": email})

	q, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQueryRow(ctx, user, q, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNoResults
	}

	return user, err
}
