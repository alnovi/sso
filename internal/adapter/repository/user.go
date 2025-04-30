package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/alnovi/sso/internal/entity"
)

const UserTable = "users"

var userFields = []string{"id", "name", "email", "password", "created_at", "updated_at", "deleted_at"}

func (r *Repository) Users(ctx context.Context, opts ...OptSelect) ([]*entity.User, error) {
	users := make([]*entity.User, 0)

	builder := r.qb.Select(userFields...).From(UserTable)
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQuery(ctx, &users, query, args...)

	return users, r.checkErr(err)
}

func (r *Repository) UsersCount(ctx context.Context, opts ...OptSelect) (int, error) {
	count := 0

	builder := r.qb.Select("COUNT (id)").From(UserTable)
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return count, err
	}

	err = r.db.QueryRow(ctx, query, args...).Scan(&count)

	return count, r.checkErr(err)
}

func (r *Repository) UserById(ctx context.Context, id string, opts ...OptSelect) (*entity.User, error) {
	user := new(entity.User)

	if err := r.checkUUID(id); err != nil {
		return nil, err
	}

	builder := r.qb.Select(userFields...).
		From(UserTable).
		Where(sq.Eq{"id": id})

	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQueryRow(ctx, user, query, args...)

	return user, r.checkErr(err)
}

func (r *Repository) UserByIds(ctx context.Context, ids []string, opts ...OptSelect) ([]*entity.User, error) {
	users := make([]*entity.User, 0)

	builder := r.qb.Select(userFields...).From(UserTable).Where(sq.Eq{"id": ids})
	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQuery(ctx, &users, query, args...)

	return users, r.checkErr(err)
}

func (r *Repository) UserByEmail(ctx context.Context, email string, opts ...OptSelect) (*entity.User, error) {
	user := new(entity.User)

	builder := r.qb.Select(userFields...).
		From(UserTable).
		Where(sq.Eq{"email": email})

	builder = r.applyOptSelect(builder, opts)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	err = r.db.ScanQueryRow(ctx, user, query, args...)

	return user, r.checkErr(err)
}

func (r *Repository) UserCreate(ctx context.Context, user *entity.User) error {
	now := time.Now()

	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}

	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	user.Id = uuid.NewString()
	user.DeletedAt = nil

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
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) UserUpdate(ctx context.Context, user *entity.User) error {
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
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) UserDelete(ctx context.Context, user *entity.User) error {
	now := time.Now()
	user.DeletedAt = &now

	builder := r.qb.Update(UserTable).
		Set("deleted_at", user.DeletedAt).
		Where(sq.Eq{"id": user.Id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}

func (r *Repository) UserDeleteForce(ctx context.Context, user *entity.User) error {
	builder := r.qb.Delete(UserTable).Where(sq.Eq{"id": user.Id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)

	return r.checkErr(err)
}
