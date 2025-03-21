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

type Option func(c *Client) error

func WithLogger(logger *slog.Logger) Option {
	return func(c *Client) error {
		c.logger = logger
		return nil
	}
}

type Client struct {
	master *pgxpool.Pool
	logger *slog.Logger
}

func NewClient(dsn string, opts ...Option) (*Client, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, err
	}

	client := &Client{master: pool}

	for _, opt := range opts {
		if err = opt(client); err != nil {
			return nil, err
		}
	}

	return client, nil
}

func (c *Client) Master() *pgxpool.Pool {
	return c.master
}

func (c *Client) DB() *sql.DB {
	return stdlib.OpenDBFromPool(c.master)
}

func (c *Client) Ping(ctx context.Context) error {
	return c.master.Ping(ctx)
}

func (c *Client) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	c.logQuery(query, args)
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, query, args...)
	}
	return c.master.Exec(ctx, query, args...)
}

func (c *Client) Query(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	c.logQuery(query, args)
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, query, args...)
	}
	return c.master.Query(ctx, query, args...)
}

func (c *Client) QueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	c.logQuery(query, args)
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, query, args...)
	}
	return c.master.QueryRow(ctx, query, args...)
}

func (c *Client) ScanQuery(ctx context.Context, dst any, query string, args ...any) error {
	rows, err := c.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	defer func() {
		rows.Close()
	}()
	return pgxscan.ScanAll(dst, rows)
}

func (c *Client) ScanQueryRow(ctx context.Context, dst any, query string, args ...any) error {
	rows, err := c.master.Query(ctx, query, args...)
	if err != nil {
		return err
	}
	return pgxscan.ScanOne(dst, rows)
}

func (c *Client) Close() error {
	c.master.Close()
	return nil
}

func (c *Client) logQuery(query string, args []any) {
	if c.logger != nil {
		c.logger.Debug(query, args)
	}
}
