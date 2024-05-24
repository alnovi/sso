package migrations

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"

	"github.com/pressly/goose/v3"
)

const (
	userPasswordCost = 14
)

func init() {
	goose.AddMigrationContext(upAddSystemUsers, downAddSystemUsers)
}

func upAddSystemUsers(ctx context.Context, tx *sql.Tx) error {
	var err error

	password, err := bcrypt.GenerateFromPassword([]byte("admin"), userPasswordCost)
	if err != nil {
		return err
	}

	query := `insert into users(id, name, email, password) values ($1, 'Admin', 'admin@example.ru', $2);`
	_, err = tx.ExecContext(ctx, query, UserAdminID(ctx), password)

	return err
}

func downAddSystemUsers(ctx context.Context, tx *sql.Tx) error {
	return nil
}
