package repository

import (
	"context"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"

	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/helper"
)

const UserTable = "users"

var userFields = []string{"id", "name", "email", "password", "created_at", "updated_at", "deleted_at"}

func (r *Repository) Users(ctx context.Context, opts ...OptSelect) ([]*entity.User, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.Users")
	defer span.End()

	users := make([]*entity.User, 0)

	builder := r.qb.Select(userFields...).From(UserTable)
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQuery(ctx, &users, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return users, nil
}

func (r *Repository) UsersCount(ctx context.Context, opts ...OptSelect) (int, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.Users")
	defer span.End()

	count := 0

	builder := r.qb.Select("COUNT (id)").From(UserTable)
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

func (r *Repository) UserById(ctx context.Context, id string, opts ...OptSelect) (*entity.User, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.UserById", helper.SpanAttr(
		attribute.String("user.id", id),
	))
	defer span.End()

	user := new(entity.User)

	if err := r.checkUUID(id); err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	builder := r.qb.Select(userFields...).
		From(UserTable).
		Where(sq.Eq{"id": id})

	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQueryRow(ctx, user, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return user, nil
}

func (r *Repository) UserByIds(ctx context.Context, ids []string, opts ...OptSelect) ([]*entity.User, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.UserByIds", helper.SpanAttr(
		attribute.String("user.ids", strings.Join(ids, ", ")),
	))
	defer span.End()

	users := make([]*entity.User, 0)

	builder := r.qb.Select(userFields...).From(UserTable).Where(sq.Eq{"id": ids})
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQuery(ctx, &users, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return users, nil
}

func (r *Repository) UserByEmail(ctx context.Context, email string, opts ...OptSelect) (*entity.User, error) {
	ctx, span := helper.SpanStart(ctx, "Repository.UserByEmail", helper.SpanAttr(
		attribute.String("user.email", email),
	))
	defer span.End()

	user := new(entity.User)

	builder := r.qb.Select(userFields...).
		From(UserTable).
		Where(sq.Eq{"email": email})

	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	err = r.checkErr(r.db.ScanQueryRow(ctx, user, query, args...))
	if err != nil {
		helper.SpanError(span, err)
		return nil, err
	}

	return user, nil
}

func (r *Repository) UserCreate(ctx context.Context, user *entity.User) error {
	ctx, span := helper.SpanStart(ctx, "Repository.UserCreate")
	defer span.End()

	now := time.Now()

	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}

	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	user.Id = uuid.NewString()
	user.DeletedAt = nil

	span.SetAttributes(attribute.String("user.id", user.Id))

	builder := r.qb.Insert(UserTable).
		Columns(userFields...).
		Values(
			user.Id,
			user.Name,
			user.Email,
			user.Password,
			user.CreatedAt,
			user.UpdatedAt,
			user.DeletedAt,
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

func (r *Repository) UserUpdate(ctx context.Context, user *entity.User) error {
	ctx, span := helper.SpanStart(ctx, "Repository.UserUpdate", helper.SpanAttr(
		attribute.String("user.id", user.Id),
	))
	defer span.End()

	user.UpdatedAt = time.Now()

	builder := r.qb.Update(UserTable).
		Set("name", user.Name).
		Set("email", user.Email).
		Set("password", user.Password).
		Set("updated_at", user.UpdatedAt).
		Set("deleted_at", user.DeletedAt).
		Where(sq.Eq{"id": user.Id})

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

func (r *Repository) UserDelete(ctx context.Context, user *entity.User) error {
	ctx, span := helper.SpanStart(ctx, "Repository.UserDelete", helper.SpanAttr(
		attribute.String("user.id", user.Id),
	))
	defer span.End()

	now := time.Now()
	user.DeletedAt = &now

	builder := r.qb.Update(UserTable).
		Set("deleted_at", user.DeletedAt).
		Where(sq.Eq{"id": user.Id})

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

func (r *Repository) UserDeleteForce(ctx context.Context, user *entity.User) error {
	ctx, span := helper.SpanStart(ctx, "Repository.UserDeleteForce", helper.SpanAttr(
		attribute.String("user.id", user.Id),
	))
	defer span.End()

	builder := r.qb.Delete(UserTable).Where(sq.Eq{"id": user.Id})

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

	return err
}
