package postgres

import (
	"context"
	"time"

	"github.com/alnovi/sso/internal/entity"
)

var tokenFields = []string{
	"id",
	"type",
	"hash",
	"user_id",
	"client_id",
	"payload",
	"not_before",
	"expiration",
	"created_at",
	"updated_at",
}

func (r *Repository) CreateToken(ctx context.Context, token *entity.Token) error {
	now := time.Now()

	token.CreatedAt = now
	token.UpdatedAt = now

	if token.Payload == nil {
		token.Payload = entity.Payload{}
	}

	return r.qb.Insert(tableTokens).
		Columns(tokenFields[1:]...).
		Suffix("RETURNING id").
		Values(
			token.Type,
			token.Hash,
			token.UserID,
			token.ClientID,
			token.Payload,
			token.NotBefore,
			token.Expiration,
			token.CreatedAt,
			token.UpdatedAt,
		).
		RunWith(r.db).
		QueryRowContext(ctx).
		Scan(&token.ID)
}
