package route

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Totus-Floreo/asperitas-on-go/internal/model"
	"github.com/Totus-Floreo/asperitas-on-go/internal/model/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type TestCase struct {
	Name      string
	Request   map[string]interface{}
	Token     string
	Response  map[string]interface{}
	AuthError error
	IsError   bool
	HTTPCode  int
}

func TestUserHandler_SignUp(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Name: "SignUp Success",
			Request: map[string]interface{}{
				"username": "user",
				"password": "password",
			},
			Token: "token",
			Response: map[string]interface{}{
				"token": "token",
			},
			AuthError: nil,
			IsError:   false,
			HTTPCode:  http.StatusCreated,
		},
		TestCase{
			Name: "SignUp User Exists Error",
			Request: map[string]interface{}{
				"username": "user",
				"password": "password",
			},
			Token: "",
			Response: map[string]interface{}{
				"error": "error",
			},
			AuthError: model.ErrUserExist,
			IsError:   true,
			HTTPCode:  http.StatusUnprocessableEntity,
		},
		TestCase{
			Name: "SignUp db error",
			Request: map[string]interface{}{
				"username": "user",
				"password": "password",
			},
			Token: "",
			Response: map[string]interface{}{
				"error": "error",
			},
			AuthError: errors.New("db error"),
			IsError:   true,
			HTTPCode:  http.StatusInternalServerError,
		},
	}

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authService := mocks.NewMockIAuthService(ctrl)

	userHandler := &UserHandler{
		Logger:      logger,
		AuthService: authService,
	}

	ts := httptest.NewServer(http.HandlerFunc(userHandler.SignUp))
	defer ts.Close()

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {

			authService.EXPECT().SignUp(gomock.Any(), gomock.Any(), gomock.Any()).Return(test.Token, test.AuthError)

			reqJSON, _ := json.Marshal(test.Request)
			r, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqJSON))
			if err != nil {
				require.NoError(t, err)
			}
			r.Header.Add("Content-Type", "application/json")

			res, err := ts.Client().Do(r)
			if err != nil {
				require.NoError(t, err)
			}
			defer res.Body.Close()

			if test.IsError {
				require.True(t, test.HTTPCode == res.StatusCode)
			} else {
				var token map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&token)
				require.NoError(t, err)
				require.Equal(t, test.Token, token["token"])
			}
		})
	}
}

func TestUserHandler_SignUpDecodeError(t *testing.T) {
	test := TestCase{
		Name: "SignUp Decoder Error",
		Request: map[string]interface{}{
			"error": "error",
		},
		Token: "",
		Response: map[string]interface{}{
			"error": "error",
		},
		AuthError: nil,
		IsError:   true,
		HTTPCode:  http.StatusBadRequest,
	}

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authService := mocks.NewMockIAuthService(ctrl)

	authService.EXPECT().SignUp(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	userHandler := &UserHandler{
		Logger:      logger,
		AuthService: authService,
	}

	ts := httptest.NewServer(http.HandlerFunc(userHandler.SignUp))
	defer ts.Close()

	reqJSON, _ := json.Marshal(test.Request)
	r, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqJSON))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/json")

	res, err := ts.Client().Do(r)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, test.HTTPCode, res.StatusCode)
}

func TestUserHandler_LogIn(t *testing.T) {
	cases := []TestCase{
		TestCase{
			Name: "LogIn Success",
			Request: map[string]interface{}{
				"username": "user",
				"password": "password",
			},
			Token: "token",
			Response: map[string]interface{}{
				"token": "token",
			},
			AuthError: nil,
			IsError:   false,
			HTTPCode:  http.StatusOK,
		},
		TestCase{
			Name: "LogIn User Invalid Credentials",
			Request: map[string]interface{}{
				"username": "user",
				"password": "password",
			},
			Token: "",
			Response: map[string]interface{}{
				"error": "error",
			},
			AuthError: model.ErrInvalidCredentials,
			IsError:   true,
			HTTPCode:  http.StatusUnauthorized,
		},
		TestCase{
			Name: "LogIn db error",
			Request: map[string]interface{}{
				"username": "user",
				"password": "password",
			},
			Token: "",
			Response: map[string]interface{}{
				"error": "error",
			},
			AuthError: errors.New("db error"),
			IsError:   true,
			HTTPCode:  http.StatusInternalServerError,
		},
	}

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authService := mocks.NewMockIAuthService(ctrl)

	userHandler := &UserHandler{
		Logger:      logger,
		AuthService: authService,
	}

	ts := httptest.NewServer(http.HandlerFunc(userHandler.LogIn))
	defer ts.Close()

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {

			authService.EXPECT().LogIn(gomock.Any(), gomock.Any(), gomock.Any()).Return(test.Token, test.AuthError)

			reqJSON, _ := json.Marshal(test.Request)
			r, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqJSON))
			if err != nil {
				require.NoError(t, err)
			}
			r.Header.Add("Content-Type", "application/json")

			res, err := ts.Client().Do(r)
			if err != nil {
				require.NoError(t, err)
			}
			defer res.Body.Close()

			if test.IsError {
				require.True(t, test.HTTPCode == res.StatusCode)
			} else {
				var token map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&token)
				require.NoError(t, err)
				require.Equal(t, test.Token, token["token"])
			}
		})
	}
}

func TestUserHandler_LogInDecodeError(t *testing.T) {
	test := TestCase{
		Name: "SignUp Decoder Error",
		Request: map[string]interface{}{
			"error": "error",
		},
		Token: "",
		Response: map[string]interface{}{
			"error": "error",
		},
		AuthError: nil,
		IsError:   true,
		HTTPCode:  http.StatusBadRequest,
	}

	zapLogger, _ := zap.NewProduction()
	defer zapLogger.Sync()
	logger := zapLogger.Sugar()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authService := mocks.NewMockIAuthService(ctrl)

	authService.EXPECT().SignUp(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	userHandler := &UserHandler{
		Logger:      logger,
		AuthService: authService,
	}

	ts := httptest.NewServer(http.HandlerFunc(userHandler.LogIn))
	defer ts.Close()

	reqJSON, _ := json.Marshal(test.Request)
	r, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(reqJSON))
	require.NoError(t, err)
	r.Header.Add("Content-Type", "application/json")

	res, err := ts.Client().Do(r)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, test.HTTPCode, res.StatusCode)
}
