package postgres

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/alnovi/sso/internal/entity"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

var clientFields = []string{
	"id",
	"name",
	"description",
	"icon",
	"color",
	"image",
	"secret",
	"home",
	"callback",
	"grant_types",
	"is_active",
	"created_at",
	"updated_at",
}

func (r *Repository) ClientByID(ctx context.Context, id string) (*entity.Client, error) {
	var err error

	if _, err = uuid.Parse(id); err != nil {
		return nil, err
	}

	client := &entity.Client{}

	err = r.qb.Select(clientFields...).
		From(tableClients).
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		QueryRowContext(ctx).
		Scan(
			&client.ID,
			&client.Name,
			&client.Description,
			&client.Icon,
			&client.Color,
			&client.Image,
			&client.Secret,
			&client.Home,
			&client.Callback,
			pq.Array(&client.GrantTypes),
			&client.IsActive,
			&client.CreatedAt,
			&client.UpdatedAt,
		)

	return client, err
}
