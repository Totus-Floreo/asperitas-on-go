package repository

import (
	"context"

	"github.com/Totus-Floreo/asperitas-on-go/pkg/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserStorage struct {
	connPool *pgxpool.Pool
}

func NewUserStorage(connPool *pgxpool.Pool) *UserStorage {
	return &UserStorage{
		connPool: connPool,
	}
}

func (s *UserStorage) GetUser(ctx context.Context, userID string) (*model.User, error) {
	conn, err := s.connPool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	tx, err := conn.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// err := tx.QueryRow(ctx, "SELECT * FROM users WHERE users.")

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &model.User{}, nil
}
