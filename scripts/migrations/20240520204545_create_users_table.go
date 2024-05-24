package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateUsersTable, downCreateUsersTable)
}

func upCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	query := `
        create table users (
            id         uuid primary key      default gen_random_uuid(),
            image      varchar,
            name       varchar      not null,
            email      varchar(100) not null,
            password   varchar      not null,
            created_at timestamptz(6) not null default now(),
            updated_at timestamptz(6) not null default now(),
            constraint users_email_unique unique (email)
        );
    `

	_, err := tx.ExecContext(ctx, query)

	return err
}

func downCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "drop table if exists users")
	return err
}
