package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ctxKeyTx struct{}

func withTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, ctxKeyTx{}, tx)
}

func txFromCtx(ctx context.Context) pgx.Tx {
	tx, _ := ctx.Value(ctxKeyTx{}).(pgx.Tx)
	return tx
}

type Transactor struct {
	pool *pgxpool.Pool
}

func NewTransactor(pool *pgxpool.Pool) *Transactor {
	return &Transactor{pool: pool}
}

func (t *Transactor) WithinTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := t.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := fn(withTx(ctx, tx)); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
