package postgres

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/alnovi/sso/internal/entity"
	"github.com/google/uuid"
)

var userFields = []string{
	"id",
	"image",
	"name",
	"email",
	"password",
	"created_at",
	"updated_at",
}

func (r *Repository) UpdateUser(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()

	_, err := r.qb.Update(tableUsers).
		Set("image", user.Image).
		Set("name", user.Name).
		Set("email", user.Email).
		Set("password", user.Password).
		Set("updated_at", user.UpdatedAt).
		Where(squirrel.Eq{"id": user.ID}).
		RunWith(r.db).
		ExecContext(ctx)

	return err
}

func (r *Repository) UserByID(ctx context.Context, id string) (*entity.User, error) {
	var err error

	if _, err = uuid.Parse(id); err != nil {
		return nil, err
	}

	user := &entity.User{}

	err = r.qb.Select(userFields...).
		From(tableUsers).
		Where(squirrel.Eq{"id": id}).
		RunWith(r.db).
		QueryRowContext(ctx).
		Scan(
			&user.ID,
			&user.Image,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

	return user, err
}

func (r *Repository) UserByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := &entity.User{}

	err := r.qb.Select(userFields...).
		From(tableUsers).
		Where(squirrel.Eq{"email": email}).
		RunWith(r.db).
		QueryRowContext(ctx).
		Scan(
			&user.ID,
			&user.Image,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)

	return user, err
}
