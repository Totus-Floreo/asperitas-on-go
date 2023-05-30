package pgx_repository

import (
	"context"
	"errors"
	"testing"

	"github.com/Totus-Floreo/asperitas-on-go/internal/model"
	"github.com/Totus-Floreo/asperitas-on-go/internal/model/mocks"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

/*
Записки для себя =)
Статья: https://medium.com/@eftal/how-to-mock-database-with-gomock-9bd0a92ffc10
Используем mockgen для генерации моков на основе интерфейсов используемых в приложении
mockgen --destination mocks/Tx.go --package=mocks  --build_flags=--mod=mod github.com/jackc/pgx/v5 Tx
mockgen --destination mocks/Row.go --package=mocks  --build_flags=--mod=mod github.com/jackc/pgx/v5 Row

Если в библиотеке нету интерфейса, то необходимо создать самому и встроить в приложение
github.com/Totus-Floreo/asperitas-on-go/internal/model/pool_interface.go
mockgen -source pool_interface.go -destination mocks/Pool.go -package=mocks
*/

type TestCase struct {
	ID      int
	Input   string
	Result  TestData
	IsError bool
	Error   error
}

type TestData struct {
	User *model.User
	Err  error
}

var ErrBegin error = errors.New("begin error")
var ErrScan error = errors.New("scan error")
var ErrCommit error = errors.New("commit error")
var ErrExec error = errors.New("exec error")

var Test = TestCase{
	ID:    0,
	Input: "1Username",
	Result: TestData{
		User: &model.User{
			ID:       "1",
			Username: "1Username",
			Password: "1Password",
		},
		Err: nil,
	},
	IsError: false,
	Error:   nil,
}

func TestGetUser_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	userStorage := NewUserStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)

	pool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().QueryRow(gomock.Any(), gomock.Any(), Test.Input).Return(mockRow)

	mockRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
		// Сохроняем в переменныую ссылку на аргумент типа строки, и потом по этой ссылке присваиваем значение
		userID := args[0].(*string)
		username := args[1].(*string)
		password := args[2].(*string)
		*userID = Test.Result.User.ID
		*username = Test.Result.User.Username
		*password = Test.Result.User.Password
		return nil
	})

	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	user, err := userStorage.GetUser(ctx, Test.Input)

	require.NoError(t, err)
	require.Equal(t, *Test.Result.User, *user)
}

func TestGetUser_UserNotFound(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	userStorage := NewUserStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)

	pool.EXPECT().Begin(ctx).Return(mockTx, nil)

	mockTx.EXPECT().QueryRow(gomock.Any(), gomock.Any(), Test.Input).Return(mockRow)

	mockRow.EXPECT().Scan(gomock.Any()).Return(pgx.ErrNoRows)

	mockTx.EXPECT().Rollback(ctx).Return(nil)

	user, err := userStorage.GetUser(ctx, Test.Input)

	require.Error(t, err)
	require.True(t, errors.Is(err, model.ErrUserNotFound))
	require.Nil(t, user)
}

func TestGetUser_BeginError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	userStorage := NewUserStorage(pool)

	pool.EXPECT().Begin(ctx).Return(nil, ErrBegin)

	user, err := userStorage.GetUser(ctx, Test.Input)

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrBegin))
	require.Nil(t, user)
}

func TestGetUser_ScanError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	userStorage := NewUserStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)

	pool.EXPECT().Begin(ctx).Return(mockTx, nil)

	mockTx.EXPECT().QueryRow(gomock.Any(), gomock.Any(), Test.Input).Return(mockRow)

	mockRow.EXPECT().Scan(gomock.Any()).Return(ErrScan)

	mockTx.EXPECT().Rollback(ctx).Return(nil)

	user, err := userStorage.GetUser(ctx, Test.Input)

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrScan))
	require.Nil(t, user)
}

func TestGetUser_CommitError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	userStorage := NewUserStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)
	mockRow := mocks.NewMockRow(ctrl)

	pool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().QueryRow(gomock.Any(), gomock.Any(), Test.Input).Return(mockRow)

	mockRow.EXPECT().Scan(gomock.Any()).DoAndReturn(func(args ...interface{}) error {
		// Сохроняем в переменныую ссылку на аргумент типа строки, и потом по этой ссылке присваиваем значение
		userID := args[0].(*string)
		username := args[1].(*string)
		password := args[2].(*string)
		*userID = Test.Result.User.ID
		*username = Test.Result.User.Username
		*password = Test.Result.User.Password
		return nil
	})

	mockTx.EXPECT().Commit(gomock.Any()).Return(ErrCommit)

	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	user, err := userStorage.GetUser(ctx, Test.Input)

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrCommit))
	require.Nil(t, user)
}

func TestAddUser_Success(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	userStorage := NewUserStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)

	pool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	mockTx.EXPECT().Commit(gomock.Any()).Return(nil)

	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	err := userStorage.AddUser(ctx, Test.Result.User)

	require.NoError(t, err)
}

func TestAddUser_BeginError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	userStorage := NewUserStorage(pool)

	pool.EXPECT().Begin(gomock.Any()).Return(nil, ErrBegin)

	err := userStorage.AddUser(ctx, Test.Result.User)

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrBegin))
}

func TestAddUser_ExecError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	userStorage := NewUserStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)

	pool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, ErrExec)

	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	err := userStorage.AddUser(ctx, Test.Result.User)

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrExec))
}

func TestAddUser_CommitError(t *testing.T) {
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pool := mocks.NewMockIPool(ctrl)
	userStorage := NewUserStorage(pool)

	mockTx := mocks.NewMockTx(ctrl)

	pool.EXPECT().Begin(gomock.Any()).Return(mockTx, nil)

	mockTx.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(pgconn.CommandTag{}, nil)

	mockTx.EXPECT().Commit(gomock.Any()).Return(ErrCommit)

	mockTx.EXPECT().Rollback(gomock.Any()).Return(nil)

	err := userStorage.AddUser(ctx, Test.Result.User)

	require.Error(t, err)
	require.True(t, errors.Is(err, ErrCommit))
}
