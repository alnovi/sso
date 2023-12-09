package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
)

func (r *Repository) CanUseClient(ctx context.Context, client entity.Client, user entity.User) error {
	canUse := client.CanUse

	err := r.qb.Select("can_use").
		From(tablePermission).
		Where(squirrel.Eq{
			"client_id": client.Id,
			"user_id":   user.Id,
		}).
		RunWith(r.connect(ctx)).
		QueryRow().
		Scan(&canUse)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if !canUse {
		return exception.AccessDenied
	}

	return nil
}
