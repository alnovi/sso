package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/alnovi/sso/internal/entity"
	"github.com/alnovi/sso/internal/exception"
	"github.com/google/uuid"
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
			&client.IsActive,
			&client.CreatedAt,
			&client.UpdatedAt,
		)

	return client, err
}

func (r *Repository) ClientByIdAndSecret(ctx context.Context, id, secret string) (*entity.Client, error) {
	client, err := r.ClientByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if !client.IsActive || client.Secret != secret {
		return nil, exception.ErrClientNotFound
	}

	return client, err
}

func (r *Repository) ClientTokenByClassAndHash(ctx context.Context, clientID, class, hash string) (*entity.Token, error) {
	token, err := r.TokenByClassAndHash(ctx, class, hash)
	if err != nil {
		return nil, err
	}

	if token.ClientID != clientID {
		return nil, fmt.Errorf("token does not belong to the client")
	}

	return token, token.IsActive()
}
