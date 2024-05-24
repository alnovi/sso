package migrations

import (
	"context"
	"database/sql"
	"time"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddTestTokens, downAddTestTokens)
}

func upAddTestTokens(ctx context.Context, tx *sql.Tx) error {
	if !EnvironmentIsTesting(ctx) {
		return nil
	}

	query := `
		insert into tokens(type, hash, user_id, client_id, expiration)
		values ('code', 'hash-code-1', $1, $2, $3);
	`

	_, err := tx.ExecContext(ctx, query,
		UserAdminID(ctx),
		ClientAdminID(ctx),
		time.Now().Add(time.Minute).Format("2006-01-02 15:04:05.000000 Z07:00"),
	)

	return err
}

func downAddTestTokens(ctx context.Context, tx *sql.Tx) error {
	return nil
}
