package postgres

import (
	"github.com/Masterminds/squirrel"

	"github.com/alnovi/sso/pkg/db/pgs"
)

type Repository struct {
	db *pgs.Client
	qb squirrel.StatementBuilderType
}

func New(db *pgs.Client) *Repository {
	return &Repository{db: db, qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
}
