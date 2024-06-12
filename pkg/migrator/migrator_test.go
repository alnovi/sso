package migrator

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name    string
		dialect string
		log     *slog.Logger
		expErr  error
	}{
		{
			name:    "Default logger",
			dialect: "postgres",
			log:     nil,
			expErr:  nil,
		},
		{
			name:    "Invalid dialect",
			dialect: "invalid",
			log:     slog.Default(),
			expErr:  errors.New("migrator can't set dialect: \"invalid\": unknown dialect"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := New(&sql.DB{}, tc.dialect, tc.log)
			assert.Equal(t, tc.expErr, err)
		})
	}
}

func TestMigrator_MigrateUp(t *testing.T) {
	m, err := New(&sql.DB{}, "postgres", nil)
	assert.NoError(t, err)
	assert.Error(t, m.MigrateUp(context.Background()))
}

func TestMigrator_MigrateDown(t *testing.T) {
	m, err := New(&sql.DB{}, "postgres", nil)
	assert.NoError(t, err)
	assert.Error(t, m.MigrateDown(context.Background()))
}
