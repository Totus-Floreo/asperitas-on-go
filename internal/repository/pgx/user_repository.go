package pgx_repository

import (
	"context"

	"github.com/Totus-Floreo/asperitas-on-go/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserStorage struct {
	connPool model.IPool
}

func NewUserStorage(connPool model.IPool) *UserStorage {
	return &UserStorage{
		connPool: connPool,
	}
}

func (s *UserStorage) GetUser(ctx context.Context, username string) (*model.User, error) {
	tx, err := s.connPool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	user := &model.User{}
	if err := tx.QueryRow(ctx, "SELECT * FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password); err != nil {
		if err == pgx.ErrNoRows {
			return nil, model.ErrUserNotFound
		} else {
			return nil, err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStorage) AddUser(ctx context.Context, user *model.User) error {
	tx, err := s.connPool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	user.ID = uuid.New().String()
	_, err = tx.Exec(ctx, "INSERT INTO users(id, username, password) VALUES ($1, $2, $3)", user.ID, user.Username, user.Password)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}
