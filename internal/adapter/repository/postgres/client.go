package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
)

var clientFields = []string{
	"id",
	"class",
	"name",
	"description",
	"logo",
	"image",
	"secret",
	"callback",
	"can_use",
	"created_at",
	"updated_at",
}

func (r *Repository) GetClientByID(ctx context.Context, id string) (*entity.Client, error) {
	result := &entity.Client{}

	err := r.qb.Select(clientFields...).
		From(tableClients).
		Where(squirrel.Eq{"id": id}).
		RunWith(r.connect(ctx)).
		QueryRow().
		Scan(
			&result.Id,
			&result.Class,
			&result.Name,
			&result.Description,
			&result.Logo,
			&result.Image,
			&result.Secret,
			&result.Callback,
			&result.CanUse,
			&result.CreatedAt,
			&result.UpdatedAt,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, exception.ClientNotFound
	}

	return result, err
}

func (r *Repository) GetClientByClass(ctx context.Context, class string) (*entity.Client, error) {
	result := &entity.Client{}

	err := r.qb.Select(clientFields...).
		From(tableClients).
		Where(squirrel.Eq{"class": class}).
		RunWith(r.connect(ctx)).
		QueryRow().
		Scan(
			&result.Id,
			&result.Class,
			&result.Name,
			&result.Description,
			&result.Logo,
			&result.Image,
			&result.Secret,
			&result.Callback,
			&result.CanUse,
			&result.CreatedAt,
			&result.UpdatedAt,
		)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, exception.ClientNotFound
	}

	return result, err
}
