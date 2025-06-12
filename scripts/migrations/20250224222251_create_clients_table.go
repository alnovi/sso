package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"

	"github.com/alnovi/sso/config"
)

func init() {
	goose.AddMigrationContext(upCreateClientsTable, downCreateClientsTable)
}

func upCreateClientsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		create table if not exists clients (
		    id          varchar(50)    primary key,
            name        varchar(50)    not null,
            secret      varchar        not null,
            callback    varchar        not null,
            icon        varchar,
            is_system   boolean        not null default false,
            created_at  timestamptz(6) not null default now(),
            updated_at  timestamptz(6) not null default now(),
		    deleted_at  timestamptz(6)
		)
	`)
	return err
}

func downCreateClientsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `drop table if exists clients;`)
	return err
}

func getConfig(ctx context.Context) *config.Config {
	return ctx.Value(config.CtxConfigKey).(*config.Config)
}
