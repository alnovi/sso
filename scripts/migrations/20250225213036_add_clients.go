package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddClients, downAddClients)
}

func upAddClients(ctx context.Context, tx *sql.Tx) error {
	cfg := getConfig(ctx)

	query := "insert into clients (id, name, icon, secret, callback, is_system) values ($1, $2, $3, $4, $5, $6)"

	args := []any{cfg.CAdmin.Id, cfg.CAdmin.Name, "/public/users.png", cfg.CAdmin.Secret, cfg.CAdmin.Callback, true}

	_, err := tx.ExecContext(ctx, query, args...)

	return err
}

func downAddClients(ctx context.Context, tx *sql.Tx) error {
	cfg := getConfig(ctx)

	_, err := tx.ExecContext(ctx, `delete from clients where id in ($1)`, cfg.CAdmin.Id)

	return err
}
