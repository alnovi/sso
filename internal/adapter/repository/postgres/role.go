package postgres

import (
	"context"
	"database/sql"
	"errors"

	sq "github.com/Masterminds/squirrel"

	"github.com/alnovi/sso/internal/entity"
)

const RoleTable = "roles"

var roleFields = []string{"client_id", "user_id", "role"}

func (r *Repository) RoleByClientAndUser(ctx context.Context, clientId, userId string) (*entity.Role, error) {
	role := new(entity.Role)

	query := r.qb.Select(roleFields...).
		From(RoleTable).
		Where(sq.Eq{"client_id": clientId, "user_id": userId})

	q, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQueryRow(ctx, role, q, args...)
	if errors.Is(err, sql.ErrNoRows) {
		role = &entity.Role{
			ClientId: clientId,
			UserId:   userId,
		}
	} else if err != nil {
		return nil, err
	}

	return role, nil
}
