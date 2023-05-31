package model

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IPool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Pool struct {
	connPool *pgxpool.Pool
}

func (pool *Pool) Begin(ctx context.Context) (pgx.Tx, error) {
	Tx, err := pool.connPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return Tx, err
}
