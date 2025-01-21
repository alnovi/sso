package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"

	"github.com/alnovi/sso/config"
	"github.com/alnovi/sso/pkg/utils"
)

func init() {
	goose.AddMigrationContext(upInsertUsers, downInsertUsers)
}

func upInsertUsers(ctx context.Context, tx *sql.Tx) error {
	cfg := getConfig(ctx)

	password, err := utils.HashPassword(cfg.Admin.Password)
	if err != nil {
		return err
	}

	query := `insert into users (id, name, email, password) values ($1, $2, $3, $4)`

	args := []any{cfg.Admin.Id, cfg.Admin.Name, cfg.Admin.Email, password}

	if cfg.App.Environment == config.AppEnvironmentTesting {
		var password string
		if password, err = utils.HashPassword(cfg.TestUser.Password); err != nil {
			return err
		}

		query = fmt.Sprintf("%s, ($5, $6, $7, $8) ", query)
		args = append(args, cfg.TestUser.Id, cfg.TestUser.Name, cfg.TestUser.Email, password)
	}

	_, err = tx.ExecContext(ctx, query, args...)

	return err
}

func downInsertUsers(ctx context.Context, tx *sql.Tx) error {
	cfg := getConfig(ctx)

	_, err := tx.ExecContext(ctx, `delete from users where email in ($1, $2)`, cfg.Admin.Email, cfg.TestUser.Email)

	return err
}
