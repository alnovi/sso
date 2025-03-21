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
	_, err := tx.ExecContext(ctx, `
		create table if not exists users (
		    id         uuid primary key        default gen_random_uuid(),
		    name       varchar(255)   not null,
		    email      varchar(255)   not null,
		    password   varchar(255)   not null,
		    created_at timestamptz(6) not null default now(),
            updated_at timestamptz(6) not null default now(),
            deleted_at timestamptz(6),
            constraint users_email_unique unique (email)
		)
	`)
	return err
}

func downCreateUsersTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `drop table if exists users`)
	return err
}
