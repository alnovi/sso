package migrator

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pressly/goose/v3"
)

type migrateKey string

const ConfigKey migrateKey = "config"

func PostgresUpFromPath(ctx context.Context, db *sql.DB, path string, log goose.Logger) error {
	goose.SetLogger(log)

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.UpContext(ctx, db, path)
	if err != nil && !errors.Is(err, goose.ErrNoMigrations) && !errors.Is(err, goose.ErrNoMigrationFiles) {
		return err
	}

	return nil
}

func PostgresResetFromPath(ctx context.Context, db *sql.DB, path string, log goose.Logger) error {
	goose.SetLogger(log)

	err := goose.SetDialect("postgres")
	if err != nil {
		return err
	}

	err = goose.ResetContext(ctx, db, path)
	if err != nil && !errors.Is(err, goose.ErrNoMigrations) && !errors.Is(err, goose.ErrNoMigrationFiles) {
		return err
	}

	return nil
}
