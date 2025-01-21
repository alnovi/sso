package migrations

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"

	"github.com/alnovi/sso/config"
)

func init() {
	goose.AddMigrationContext(upInsertClients, downInsertClients)
}

func upInsertClients(ctx context.Context, tx *sql.Tx) error {
	cfg := getConfig(ctx)

	query := `insert into clients (id, name, secret, host) values ($1, $2, $3, $4)`

	args := []any{cfg.Client.Id, cfg.Client.Name, cfg.Client.Secret, cfg.Client.Host}

	if cfg.App.Environment == config.AppEnvironmentTesting {
		query = fmt.Sprintf("%s, ($5, $6, $7, $8) ", query)
		args = append(args, cfg.TestClient.Id, cfg.TestClient.Name, cfg.TestClient.Secret, cfg.TestClient.Host)
	}

	_, err := tx.ExecContext(ctx, query, args...)

	return err
}

func downInsertClients(ctx context.Context, tx *sql.Tx) error {
	cfg := getConfig(ctx)

	_, err := tx.ExecContext(ctx, `delete from clients where id in ($1, $2)`, cfg.Client.Id, cfg.TestClient.Id)

	return err
}
