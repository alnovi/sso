package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateClientsTable, downCreateClientsTable)
}

func upCreateClientsTable(ctx context.Context, tx *sql.Tx) error {
	query := `
        create table clients (
            id          uuid primary key        default gen_random_uuid(),
            name        varchar(50)    not null,
            description varchar,
            icon        varchar,
            color       varchar,
            image       varchar,
            secret      varchar        not null,
            home        varchar        not null,
            callback    varchar        not null,
            grant_types varchar[]      not null default array ['code']::varchar[],
            is_active   boolean        not null default true,
            created_at  timestamptz(6) not null default now(),
            updated_at  timestamptz(6) not null default now()
        );
    `

	_, err := tx.ExecContext(ctx, query)

	return err
}

func downCreateClientsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "drop table if exists clients")
	return err
}
