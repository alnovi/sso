package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddRoles, downAddRoles)
}

func upAddRoles(ctx context.Context, tx *sql.Tx) error {
	cfg := getConfig(ctx)

	query := `insert into roles (client_id, user_id, role) values ($1, $2, $3)`

	args := []any{cfg.CAdmin.Id, cfg.UAdmin.Id, "admin"}

	_, err := tx.ExecContext(ctx, query, args...)

	return err
}

func downAddRoles(ctx context.Context, tx *sql.Tx) error {
	cfg := getConfig(ctx)

	_, err := tx.ExecContext(ctx, `delete from roles where client_id = $1 and user_id = $2`, cfg.CAdmin.Id, cfg.UAdmin.Id)

	return err
}
