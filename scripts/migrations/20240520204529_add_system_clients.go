package migrations

import (
	"context"
	"database/sql"

	"github.com/alnovi/sso/pkg/rand"
	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddSystemClients, downAddSystemClients)
}

const ClientSecretCost = 50

func upAddSystemClients(ctx context.Context, tx *sql.Tx) error {
	query := `
		insert into clients(id, name, icon, color, secret, home, callback)
		values ($1, 'Администрирование', 'fab fa-gripfire', '#d9342b', $2, '/admin', '/admin/callback'),
			   ($3, 'Единый вход', 'fa fa-shield-halved', '#8b5cf6', $4, '/profile', '/profile/callback');
	`

	_, err := tx.ExecContext(ctx, query,
		ClientAdminID(ctx),
		rand.Base62(ClientSecretCost),
		ClientProfileID(ctx),
		rand.Base62(ClientSecretCost),
	)

	return err
}

func downAddSystemClients(ctx context.Context, tx *sql.Tx) error {
	return nil
}
