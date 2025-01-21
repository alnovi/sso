package postgres

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/alnovi/sso/internal/adapter/repository"
	"github.com/alnovi/sso/internal/entity"
)

const ClientTable = "clients"

var clientFields = []string{"id", "name", "secret", "host", "icon", "color", "image", "is_active", "created_at", "updated_at"}

func (r *Repository) ClientById(ctx context.Context, id string) (*entity.Client, error) {
	client := new(entity.Client)

	query := r.qb.Select(clientFields...).
		From(ClientTable).
		Where(sq.Eq{"id": id})

	q, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQueryRow(ctx, client, q, args...)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, repository.ErrNoResults
	}

	return client, err
}
