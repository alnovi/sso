package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/helper"
)

const RoleTable = "roles"

var roleFields = []string{"client_id", "user_id", "role"}

func (r *Repository) Role(ctx context.Context, clientId, userId string) (*entity.Role, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.Role", helper.SpanAttr(
		attribute.String("client.id", clientId),
		attribute.String("user.id", userId),
	))
	defer span.End()

	role := new(entity.Role)

	if err := r.checkUUID(userId); err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	builder := r.qb.Select(roleFields...).
		From(RoleTable).
		Where(sq.Eq{"client_id": clientId, "user_id": userId})

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQueryRow(ctx, role, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return role, nil
}

func (r *Repository) RoleByUserId(ctx context.Context, userId string, opts ...OptSelect) ([]*entity.Role, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.RoleByUserId", helper.SpanAttr(
		attribute.String("user.id", userId),
	))
	defer span.End()

	roles := make([]*entity.Role, 0)

	if err := r.checkUUID(userId); err != nil {
		helper.SpanError(span, err)
		return roles, nil //nolint:nilerr
	}

	builder := r.qb.Select(roleFields...).
		From(RoleTable).
		Where(sq.Eq{"user_id": userId})

	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQuery(ctx, &roles, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return roles, nil
}

func (r *Repository) RoleUpdate(ctx context.Context, role *entity.Role) error {
	ctx, span := helper.SpanStart(ctx, "Repository.RoleUpdate", helper.SpanAttr(
		attribute.String("role.role", role.Role),
		attribute.String("role.user.id", role.UserId),
		attribute.String("role.client.id", role.ClientId),
	))
	defer span.End()

	builder := r.qb.Insert(RoleTable).
		Columns(roleFields...).
		Values(role.ClientId, role.UserId, role.Role).
		Suffix(`ON CONFLICT (client_id, user_id) DO UPDATE SET role = EXCLUDED.role`)

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

func (r *Repository) RoleDelete(ctx context.Context, clientId, userId string) error {
	ctx, span := helper.SpanStart(ctx, "Repository.RoleDelete", helper.SpanAttr(
		attribute.String("role.user.id", userId),
		attribute.String("role.client.id", clientId),
	))
	defer span.End()

	builder := r.qb.Delete(RoleTable).Where(sq.Eq{"client_id": clientId, "user_id": userId})

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
