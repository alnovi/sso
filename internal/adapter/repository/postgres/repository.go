package postgres

import (
	"github.com/Masterminds/squirrel"

	"github.com/alnovi/sso/pkg/client/postgres"
)

type Repository struct {
	db *postgres.Client
	qb squirrel.StatementBuilderType
}

func New(db *postgres.Client) *Repository {
	return &Repository{db: db, qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
}
