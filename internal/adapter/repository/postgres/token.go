package postgres

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/alnovi/sso/internal/entity"
)

var tokenFields = []string{
	"id",
	"class",
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
			token.Class,
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

func (r *Repository) DeleteTokenById(ctx context.Context, id string) error {
	_, err := r.qb.Delete(tableTokens).
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		ExecContext(ctx)

	return err
}

func (r *Repository) TokenByClassAndHash(ctx context.Context, class, hash string) (*entity.Token, error) {
	token := &entity.Token{}

	err := r.qb.Select(tokenFields...).
		From(tableTokens).
		Where(sq.Eq{"class": class, "hash": hash}).
		RunWith(r.db).
		QueryRowContext(ctx).
		Scan(
			&token.ID,
			&token.Class,
			&token.Hash,
			&token.UserID,
			&token.ClientID,
			&token.Payload,
			&token.NotBefore,
			&token.Expiration,
			&token.CreatedAt,
			&token.UpdatedAt,
		)

	return token, err
}
