package migrations

import (
	"context"
	"database/sql"

	"github.com/alnovi/gomon/utils"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddUsers, downAddUsers)
}

func upAddUsers(ctx context.Context, tx *sql.Tx) error {
	cfg := getConfig(ctx)

	password, _ := utils.HashPassword(cfg.UAdmin.Password)

	query := `insert into users (id, name, email, password) values ($1, $2, $3, $4)`

	args := []any{cfg.UAdmin.Id, cfg.UAdmin.Name, cfg.UAdmin.Email, password}

	_, err := tx.ExecContext(ctx, query, args...)

	return err
}

func downAddUsers(ctx context.Context, tx *sql.Tx) error {
	cfg := getConfig(ctx)

	_, err := tx.ExecContext(ctx, `delete from users where email in ($1)`, cfg.UAdmin.Email)

	return err
}
