package redis_repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenRepository struct {
	rdb *redis.Client
}

func NewTokenRepository(rdb *redis.Client) *TokenRepository {
	return &TokenRepository{
		rdb: rdb,
	}
}

func (r *TokenRepository) GetToken(ctx context.Context, userID string) (string, error) {
	val, err := r.rdb.Get(ctx, userID).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (r *TokenRepository) SetToken(ctx context.Context, userID string, token string) error {
	err := r.rdb.SetNX(ctx, userID, token, time.Hour*24*7).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *TokenRepository) DeleteToken(ctx context.Context, userID string) error {
	err := r.rdb.Del(ctx, userID).Err()
	if err != nil {
		return err
	}
	return nil
}
