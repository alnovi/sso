package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateTokensTable, downCreateTokensTable)
}

func upCreateTokensTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
		create  table if not exists tokens (
    		id         uuid primary key default gen_random_uuid(),
    		class      varchar(30) not null,
    		hash       varchar(500) not null,
    		session_id uuid,
    		user_id    uuid,
    		client_id  varchar(50),
    		payload    jsonb        not null,
            not_before timestamptz(6) not null default now(),
            expiration timestamptz(6) not null default now(),
            created_at timestamptz(6) not null default now(),
            updated_at timestamptz(6) not null default now(),
            constraint tokens_class_hash_unique unique (class, hash),
            constraint tokens_session_fk foreign key (session_id) references sessions (id) on delete cascade on update cascade,
            constraint tokens_user_fk foreign key (user_id) references users (id) on delete cascade on update cascade,
            constraint tokens_client_fk foreign key (client_id) references clients (id) on delete cascade on update cascade
		)
	`)
	return err
}

func downCreateTokensTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `drop table if exists tokens;`)
	return err
}
