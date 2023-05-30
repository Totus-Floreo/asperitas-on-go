package redis_repository

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

type TestCase struct {
	Name   string
	UserID string
	Token  string
	Error  error
}

var Test []TestCase = []TestCase{
	TestCase{
		Name:   "Success",
		UserID: "ValidUser",
		Token:  "ValidToken",
		Error:  nil,
	},
}

func TestGetToken_Success(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	db := NewTokenRepository(client)

	mock.ExpectGet(Test[0].UserID).SetVal(Test[0].Token)

	token, err := db.GetToken(ctx, Test[0].UserID)

	require.NoError(t, err)
	require.Equal(t, token, Test[0].Token)
}

func TestGetToken_Fail(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	db := NewTokenRepository(client)

	mock.ExpectGet(Test[0].UserID).RedisNil()

	token, err := db.GetToken(ctx, Test[0].UserID)

	require.Error(t, err)
	require.True(t, errors.Is(err, redis.Nil))
	require.Empty(t, token)
}

func TestSetToken_Success(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	db := NewTokenRepository(client)

	mock.ExpectSetNX(Test[0].UserID, Test[0].Token, time.Hour*24*7).SetVal(true)

	err := db.SetToken(ctx, Test[0].UserID, Test[0].Token)

	require.NoError(t, err)
}

func TestSetToken_Fail(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	db := NewTokenRepository(client)

	mock.ExpectSetNX(Test[0].UserID, Test[0].Token, time.Hour*24*7).SetErr(redis.TxFailedErr)

	err := db.SetToken(ctx, Test[0].UserID, Test[0].Token)

	require.Error(t, err)
	require.True(t, errors.Is(err, redis.TxFailedErr))
}

func TestDeleteToken_Success(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	db := NewTokenRepository(client)

	mock.ExpectDel(Test[0].UserID).SetVal(1)

	err := db.DeleteToken(ctx, Test[0].UserID)

	require.NoError(t, err)
}

func TestDeleteToken_Fail(t *testing.T) {
	ctx := context.Background()

	client, mock := redismock.NewClientMock()
	db := NewTokenRepository(client)

	mock.ExpectDel(Test[0].UserID).SetErr(redis.TxFailedErr)

	err := db.DeleteToken(ctx, Test[0].UserID)

	require.Error(t, err)
	require.True(t, errors.Is(err, redis.TxFailedErr))
}
