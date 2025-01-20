package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const txKey key = "tx"

type Transaction struct {
	db *pgxpool.Pool
}

func NewTransaction(db *pgxpool.Pool) *Transaction {
	return &Transaction{db: db}
}

func (t *Transaction) ReadCommitted(ctx context.Context, fn func(ctx context.Context) error) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return t.transaction(ctx, txOpts, fn)
}

func (t *Transaction) transaction(ctx context.Context, opts pgx.TxOptions, fn func(ctx context.Context) error) error {
	_, ok := ctx.Value(txKey).(*pgx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err := t.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("can't begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	err = fn(context.WithValue(ctx, txKey, tx))
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("can't commit transaction: %w", err)
	}

	return nil
}
