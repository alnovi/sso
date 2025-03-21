package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upCreateRolesTable, downCreateRolesTable)
}

func upCreateRolesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
        create table if not exists roles (
            client_id   varchar(50) references clients (id) on delete cascade on update cascade,
            user_id     uuid        references users (id)   on delete cascade on update cascade,
            role        varchar(50) not null default 'user'
        );
		create unique index roles_client_id_user_id_unique on roles (client_id, user_id);
    `)
	return err
}

func downCreateRolesTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `drop table if exists roles;`)
	return err
}
