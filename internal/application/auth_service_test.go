package application

import (
	"context"
	"errors"
	"testing"

	"github.com/Totus-Floreo/asperitas-on-go/internal/model"
	"github.com/Totus-Floreo/asperitas-on-go/internal/model/mocks"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Name  string
	User  *model.User
	Token string
	Error error
}

func TestAuthServiceSignUp(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Name: "AuthService: SignUp Success",
			User: &model.User{
				ID:       "id0",
				Username: "User0",
				Password: "Password0",
			},
			Token: "token0",
			Error: nil,
		},
		TestCase{
			Name: "AuthService: SignUp Error User Exist",
			User: &model.User{
				ID:       "id1",
				Username: "User1",
				Password: "Password1",
			},
			Token: "token1",
			Error: model.ErrUserExist,
		},
		TestCase{
			Name: "AuthService: SignUp AddUser Error",
			User: &model.User{
				ID:       "id2",
				Username: "User2",
				Password: "Password2",
			},
			Token: "token2",
			Error: errors.New("db error"),
		},
		TestCase{
			Name: "AuthService: SignUp GenerateToken Error",
			User: &model.User{
				ID:       "id3",
				Username: "User3",
				Password: "Password3",
			},
			Token: "token3",
			Error: model.ErrInvalidToken,
		},
		TestCase{
			Name: "AuthService: SignUp SetToken Error",
			User: &model.User{
				ID:       "id4",
				Username: "User4",
				Password: "Password4",
			},
			Token: "token4",
			Error: redis.TxFailedErr,
		},
		TestCase{
			Name: "AuthService: LogIn Success",
			User: &model.User{
				ID:       "id5",
				Username: "User5",
				Password: "Password5",
			},
			Token: "token5",
			Error: nil,
		},
		TestCase{
			Name: "AuthService: LogIn Success but expired token",
			User: &model.User{
				ID:       "id6",
				Username: "User6",
				Password: "Password6",
			},
			Token: "token6",
			Error: nil,
		},
		TestCase{
			Name: "AuthService: LogIn GetUser Error User Not Found",
			User: &model.User{
				ID:       "id7",
				Username: "User7",
				Password: "Password7",
			},
			Token: "token7",
			Error: model.ErrInvalidCredentials,
		},
		TestCase{
			Name: "AuthService: LogIn GetUser Error Invalid Credentials",
			User: &model.User{
				ID:       "id8",
				Username: "User8",
				Password: "Password8",
			},
			Token: "token8",
			Error: model.ErrInvalidCredentials,
		},
		TestCase{
			Name: "AuthService: LogIn GetToken Error No Token",
			User: &model.User{
				ID:       "id9",
				Username: "User9",
				Password: "Password9",
			},
			Token: "token9",
			Error: errors.New("generate error"),
		},
		TestCase{
			Name: "AuthService: LogIn GenerateToken Error",
			User: &model.User{
				ID:       "id10",
				Username: "User10",
				Password: "Password10",
			},
			Token: "token10",
			Error: redis.TxFailedErr,
		},
		TestCase{
			Name: "AuthService: LogIn GetToken Error",
			User: &model.User{
				ID:       "id11",
				Username: "User11",
				Password: "Password11",
			},
			Token: "token11",
			Error: redis.TxFailedErr,
		},
	}

	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userStorage := mocks.NewMockIUserStorage(ctrl)
	tokenStorage := mocks.NewMockITokenStorage(ctrl)
	jwtService := mocks.NewMockIJWTService(ctrl)

	authService := NewAuthService(userStorage, tokenStorage, jwtService)

	t.Run(cases[0].Name, func(t *testing.T) {
		userStorage.EXPECT().GetUser(gomock.Any(), cases[0].User.Username).Return(nil, model.ErrUserNotFound)

		userStorage.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(nil)

		jwtService.EXPECT().GenerateToken(gomock.Any()).Return(cases[0].Token, nil)

		tokenStorage.EXPECT().SetToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

		token, err := authService.SignUp(ctx, cases[0].User.Username, cases[0].User.Password)

		assert.NoError(t, err)
		assert.Equal(t, cases[0].Token, token)
	})

	t.Run(cases[1].Name, func(t *testing.T) {
		userStorage.EXPECT().GetUser(gomock.Any(), cases[1].User.Username).Return(cases[1].User, nil)

		token, err := authService.SignUp(ctx, cases[1].User.Username, cases[1].User.Password)

		assert.Error(t, err)
		assert.True(t, errors.Is(cases[1].Error, err))
		assert.Empty(t, token)
	})

	t.Run(cases[2].Name, func(t *testing.T) {
		userStorage.EXPECT().GetUser(gomock.Any(), cases[2].User.Username).Return(nil, model.ErrUserNotFound)

		userStorage.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(cases[2].Error)

		token, err := authService.SignUp(ctx, cases[2].User.Username, cases[2].User.Password)

		assert.Error(t, err)
		assert.True(t, errors.Is(cases[2].Error, err))
		assert.Empty(t, token)
	})

	t.Run(cases[3].Name, func(t *testing.T) {
		userStorage.EXPECT().GetUser(gomock.Any(), cases[3].User.Username).Return(nil, model.ErrUserNotFound)

		userStorage.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(nil)

		jwtService.EXPECT().GenerateToken(gomock.Any()).Return(gomock.Nil().String(), model.ErrInvalidToken)

		token, err := authService.SignUp(ctx, cases[3].User.Username, cases[3].User.Password)

		assert.Error(t, err)
		assert.True(t, errors.Is(cases[3].Error, err))
		assert.Empty(t, token)
	})

	t.Run(cases[4].Name, func(t *testing.T) {
		userStorage.EXPECT().GetUser(gomock.Any(), cases[4].User.Username).Return(nil, model.ErrUserNotFound)

		userStorage.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(nil)

		jwtService.EXPECT().GenerateToken(gomock.Any()).Return(cases[4].Token, nil)

		tokenStorage.EXPECT().SetToken(gomock.Any(), gomock.Any(), gomock.Any()).Return(cases[4].Error)

		token, err := authService.SignUp(ctx, cases[4].User.Username, cases[4].User.Password)

		assert.Error(t, err)
		assert.True(t, errors.Is(cases[4].Error, err))
		assert.Empty(t, token)
	})

	//LogIn

	t.Run(cases[5].Name, func(t *testing.T) {
		userStorage.EXPECT().GetUser(gomock.Any(), cases[5].User.Username).Return(cases[5].User, nil)

		tokenStorage.EXPECT().GetToken(gomock.Any(), cases[5].User.ID).Return(cases[5].Token, nil)

		token, err := authService.LogIn(ctx, cases[5].User.Username, cases[5].User.Password)

		assert.NoError(t, err)
		assert.Equal(t, cases[5].Token, token)
	})

	t.Run(cases[6].Name, func(t *testing.T) {
		userStorage.EXPECT().GetUser(gomock.Any(), cases[6].User.Username).Return(cases[6].User, nil)

		tokenStorage.EXPECT().GetToken(gomock.Any(), cases[6].User.ID).Return(gomock.Nil().String(), redis.Nil)

		jwtService.EXPECT().GenerateToken(cases[6].User).Return(cases[6].Token, nil)

		tokenStorage.EXPECT().SetToken(gomock.Any(), cases[6].User.ID, cases[6].Token).Return(nil)

		token, err := authService.LogIn(ctx, cases[6].User.Username, cases[6].User.Password)

		assert.NoError(t, err)
		assert.Equal(t, cases[6].Token, token)
	})

	t.Run(cases[7].Name, func(t *testing.T) {
		userStorage.EXPECT().GetUser(gomock.Any(), cases[7].User.Username).Return(nil, model.ErrInvalidCredentials)

		token, err := authService.LogIn(ctx, cases[7].User.Username, cases[7].User.Password)

		assert.Error(t, err)
		assert.True(t, errors.Is(cases[7].Error, err))
		assert.Empty(t, token)
	})

	t.Run(cases[8].Name, func(t *testing.T) {
		userStorage.EXPECT().GetUser(gomock.Any(), cases[8].User.Username).Return(cases[8].User, nil)

		token, err := authService.LogIn(ctx, cases[8].User.Username, cases[8].User.Password+"missclick")

		assert.Error(t, err)
		assert.True(t, errors.Is(cases[8].Error, err))
		assert.Empty(t, token)
	})

	t.Run(cases[9].Name, func(t *testing.T) {
		userStorage.EXPECT().GetUser(gomock.Any(), cases[9].User.Username).Return(cases[9].User, nil)

		tokenStorage.EXPECT().GetToken(gomock.Any(), cases[9].User.ID).Return(gomock.Nil().String(), redis.Nil)

		jwtService.EXPECT().GenerateToken(cases[9].User).Return(gomock.Nil().String(), cases[9].Error)

		token, err := authService.LogIn(ctx, cases[9].User.Username, cases[9].User.Password)

		assert.Error(t, err)
		assert.True(t, errors.Is(cases[9].Error, err))
		assert.Empty(t, token)
	})

	t.Run(cases[10].Name, func(t *testing.T) {
		userStorage.EXPECT().GetUser(gomock.Any(), cases[10].User.Username).Return(cases[10].User, nil)

		tokenStorage.EXPECT().GetToken(gomock.Any(), cases[10].User.ID).Return(gomock.Nil().String(), redis.Nil)

		jwtService.EXPECT().GenerateToken(cases[10].User).Return(cases[10].Token, nil)

		tokenStorage.EXPECT().SetToken(gomock.Any(), cases[10].User.ID, cases[10].Token).Return(redis.TxFailedErr)

		token, err := authService.LogIn(ctx, cases[10].User.Username, cases[10].User.Password)

		assert.Error(t, err)
		assert.True(t, errors.Is(cases[10].Error, err))
		assert.Empty(t, token)
	})

	t.Run(cases[11].Name, func(t *testing.T) {
		userStorage.EXPECT().GetUser(gomock.Any(), cases[11].User.Username).Return(cases[11].User, nil)

		tokenStorage.EXPECT().GetToken(gomock.Any(), cases[11].User.ID).Return(gomock.Nil().String(), redis.TxFailedErr)

		token, err := authService.LogIn(ctx, cases[11].User.Username, cases[11].User.Password)

		assert.Error(t, err)
		assert.True(t, errors.Is(cases[11].Error, err))
		assert.Empty(t, token)
	})
}
