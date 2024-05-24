package migrator

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db *sql.DB
}

func New(db *sql.DB, dialect string, log *slog.Logger) (*Migrator, error) {
	if log == nil {
		log = slog.Default()
	}

	goose.SetLogger(NewLogger(log))
	err := goose.SetDialect(dialect)
	if err != nil {
		return nil, fmt.Errorf("migrator can't set dialect: %s", err)
	}

	return &Migrator{
		db: db,
	}, nil
}

func (m *Migrator) MigrateUp(ctx context.Context) error {
	if err := goose.UpContext(ctx, m.db, "."); err != nil {
		return fmt.Errorf("migrator fail up migrate: %s", err)
	}
	return nil
}

func (m *Migrator) MigrateDown(ctx context.Context) error {
	if err := goose.DownToContext(ctx, m.db, ".", 0); err != nil {
		return fmt.Errorf("migrator fail down migrate: %s", err)
	}
	return nil
}
