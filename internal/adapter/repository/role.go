package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"

	"github.com/alnovi/sso/internal/entity"
)

const RoleTable = "roles"

var roleFields = []string{"client_id", "user_id", "role"}

func (r *Repository) Role(ctx context.Context, clientId, userId string) (*entity.Role, error) {
	role := new(entity.Role)

	if err := r.checkUUID(userId); err != nil {
		return nil, err
	}

	builder := r.qb.Select(roleFields...).
		From(RoleTable).
		Where(sq.Eq{"client_id": clientId, "user_id": userId})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQueryRow(ctx, role, query, args...)

	return role, r.checkErr(err)
}

func (r *Repository) RoleByUserId(ctx context.Context, userId string, opts ...OptSelect) ([]*entity.Role, error) {
	roles := make([]*entity.Role, 0)

	if err := r.checkUUID(userId); err != nil {
		return roles, nil //nolint:nilerr
	}

	builder := r.qb.Select(roleFields...).
		From(RoleTable).
		Where(sq.Eq{"user_id": userId})

	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQuery(ctx, &roles, query, args...)

	return roles, r.checkErr(err)
}

func (r *Repository) RoleUpdate(ctx context.Context, role *entity.Role) error {
	builder := r.qb.Insert(RoleTable).
		Columns(roleFields...).
		Values(role.ClientId, role.UserId, role.Role).
		Suffix(`ON CONFLICT (client_id, user_id) DO UPDATE SET role = EXCLUDED.role`)

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) RoleDelete(ctx context.Context, clientId, userId string) error {
	builder := r.qb.Delete(RoleTable).Where(sq.Eq{"client_id": clientId, "user_id": userId})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}
