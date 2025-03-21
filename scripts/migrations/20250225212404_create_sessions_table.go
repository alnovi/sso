package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateSessionsTable, downCreateSessionsTable)
}

func upCreateSessionsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		create  table if not exists sessions (
    		id         uuid primary key default gen_random_uuid(),
    		user_id    uuid not null,
    		ip         varchar(50),
    		agent      varchar(250),
            created_at timestamptz(6) not null default now(),
            updated_at timestamptz(6) not null default now(),
            constraint sessions_user_fk foreign key (user_id) references users (id) on delete cascade on update cascade
		)
	`)
	return err
}

func downCreateSessionsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `drop table if exists sessions;`)
	return err
}
