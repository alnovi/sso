package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
)

var userFields = []string{
	"id",
	"image",
	"name",
	"login",
	"email",
	"password",
	"created_at",
	"updated_at",
}

func (r *Repository) GetUserById(ctx context.Context, id string) (*entity.User, error) {
	result := &entity.User{}

	err := r.qb.Select(userFields...).
		From(tableUsers).
		Where(squirrel.Eq{"id": id}).
		RunWith(r.connect(ctx)).
		QueryRow().
		Scan(
			&result.Id,
			&result.Image,
			&result.Name,
			&result.Login,
			&result.Email,
			&result.Password,
			&result.CreatedAt,
			&result.UpdatedAt,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, exception.UserNotFound
	}

	return result, err
}

func (r *Repository) GetUserByLoginOrEmail(ctx context.Context, login string) (*entity.User, error) {
	result := &entity.User{}

	err := r.qb.Select(userFields...).
		From(tableUsers).
		Where(squirrel.Or{
			squirrel.Eq{"login": login},
			squirrel.Eq{"email": login},
		}).
		RunWith(r.connect(ctx)).
		QueryRow().
		Scan(
			&result.Id,
			&result.Image,
			&result.Name,
			&result.Login,
			&result.Email,
			&result.Password,
			&result.CreatedAt,
			&result.UpdatedAt,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, exception.UserNotFound
	}

	return result, err
}
