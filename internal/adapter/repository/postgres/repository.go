package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/Masterminds/squirrel"
	"github.com/alnovi/sso/pkg/migrator"
	_ "github.com/alnovi/sso/scripts/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	tableClients = "clients"
	tableUsers   = "users"
	tableTokens  = "tokens"
)

var (
	mode = map[bool]string{true: "enable", false: "disable"}
)

type Repository struct {
	db *sql.DB
	qb squirrel.StatementBuilderType
}

func New(host, port, database, user, password string, ssl bool) (*Repository, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&TimeZone=UTC", user, password, host, port, database, mode[ssl])

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %s", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("can't ping database: %s", err)
	}

	return &Repository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

func (r *Repository) MigrateUp(ctx context.Context, log *slog.Logger) error {
	m, err := migrator.New(r.db, "postgres", log)
	if err != nil {
		return fmt.Errorf("postgres repository fail up migrate: %s", err)
	}

	if err = m.MigrateUp(ctx); err != nil {
		return fmt.Errorf("postgres repository fail up migrate: %s", err)
	}

	return nil
}

func (r *Repository) MigrateDown(ctx context.Context, log *slog.Logger) error {
	m, err := migrator.New(r.db, "postgres", log)
	if err != nil {
		return fmt.Errorf("postgres repository fail down migrate: %s", err)
	}

	if err = m.MigrateDown(ctx); err != nil {
		return fmt.Errorf("postgres repository fail down migrate: %s", err)
	}

	return nil
}

func (r *Repository) Close(_ context.Context) error {
	err := r.db.Close()
	if err != nil {
		return fmt.Errorf("postgres repository can't close connection: %s", err)
	}
	return nil
}
