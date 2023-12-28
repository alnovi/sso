package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/alnovi/sso/scripts"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type ctxKey string

const (
	keyTx           ctxKey = "tx"
	tableClients    string = "clients"
	tablePermission string = "permissions"
	tableTokens     string = "tokens"
	tableUsers      string = "users"
)

type Config struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
	SSLMode  bool
}

type Repository struct {
	db *sql.DB
	qb squirrel.StatementBuilderType
}

func NewRepository(cfg Config) (*Repository, error) {
	mode := map[bool]string{
		true:  "enable",
		false: "disable",
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database, mode[cfg.SSLMode])

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &Repository{
		db: db,
		qb: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}, nil
}

func (r *Repository) MigrateUp() error {
	goose.SetBaseFS(scripts.MigrateSchema)

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	return goose.Up(r.db, "migrations")
}

func (r *Repository) Close() error {
	return r.db.Close()
}

func (r *Repository) Begin(ctx context.Context) (context.Context, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	txCtx := context.WithValue(ctx, keyTx, tx)

	return txCtx, nil
}

func (r *Repository) Commit(ctx context.Context) error {
	if tx, ok := ctx.Value(keyTx).(*sql.Tx); ok {
		return tx.Commit()
	}
	return nil
}

func (r *Repository) Rollback(ctx context.Context) error {
	if tx, ok := ctx.Value(keyTx).(*sql.Tx); ok {
		return tx.Rollback()
	}
	return nil
}

func (r *Repository) connect(ctx context.Context) squirrel.BaseRunner {
	if tx, ok := ctx.Value(keyTx).(*sql.Tx); ok {
		return tx
	}
	return r.db
}
