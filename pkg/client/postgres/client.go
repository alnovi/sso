package postgres

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type key string

type TransactionFn func(ctx context.Context) error

type Client struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewClient(dsn string) (*Client, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	return &Client{db: pool}, nil
}

func (c *Client) SetLogger(logger *slog.Logger) {
	c.log = logger
}

func (c *Client) DB() *pgxpool.Pool {
	return c.db
}

func (c *Client) SqlDB() *sql.DB {
	return stdlib.OpenDBFromPool(c.db)
}

func (c *Client) Ping(ctx context.Context) error {
	return c.db.Ping(ctx)
}

func (c *Client) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	c.logQuery(query, args)
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, query, args...)
	}
	return c.db.Exec(ctx, query, args...)
}

func (c *Client) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	c.logQuery(query, args)
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, query, args...)
	}
	return c.db.Query(ctx, query, args...)
}

func (c *Client) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	c.logQuery(query, args)
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, query, args...)
	}
	return c.db.QueryRow(ctx, query, args...)
}

func (c *Client) ScanQuery(ctx context.Context, dst any, query string, args ...interface{}) error {
	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanAll(dst, rows)
}

func (c *Client) ScanQueryRow(ctx context.Context, dst any, query string, args ...interface{}) error {
	rows, err := c.db.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanOne(dst, rows)
}

func (c *Client) Close() error {
	c.db.Close()
	return nil
}

func (c *Client) logQuery(query string, args []interface{}) {
	if c.log != nil {
		c.log.Debug(query, args...)
	}
}
