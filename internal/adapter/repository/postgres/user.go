package postgres

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/alnovi/sso/internal/entity"
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
