package application

import (
	"errors"
	"testing"
	"time"

	"github.com/Totus-Floreo/asperitas-on-go/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

type TestCaseGen struct {
	Name    string
	Secret  string
	Method  jwt.SigningMethod
	Input   *model.User
	Result  string
	IsError bool
	Error   error
}

type TestCaseVerify struct {
	Name    string
	Secret  string
	Method  jwt.SigningMethod
	Input   string
	Result  *model.Author
	IsError bool
	Error   error
}

type FakeTimeController struct {
	fixedTime time.Time
}

func (t *FakeTimeController) Now() time.Time {
	return t.fixedTime
}

func TestJWTServiceGenerateToken(t *testing.T) {
	testCases := []TestCaseGen{
		TestCaseGen{
			Name:   "Success",
			Secret: "testkey",
			Method: jwt.SigningMethodHS256,
			Input: &model.User{
				ID:       "test",
				Username: "TestUser",
				Password: "12345678",
			},
			// https://jwt.io/
			Result:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoiVGVzdFVzZXIiLCJpZCI6InRlc3QifSwiZXhwIjoxMjU4NDk4ODAwLCJpYXQiOjEyNTc4OTQwMDB9.932YbwTvx-9LLM4eBCpDYt69qP9SmuSWPB7vcLamsXg",
			IsError: false,
			Error:   nil,
		},
		TestCaseGen{
			Name:   "Broken Method",
			Secret: "testkey",
			Method: jwt.SigningMethodNone,
			Input: &model.User{
				ID:       "test",
				Username: "TestUser",
				Password: "12345678",
			},
			// https://jwt.io/
			Result:  "d",
			IsError: true,
			Error:   nil,
		},
	}

	mock := new(FakeTimeController)
	mock.fixedTime = time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			jwtService := NewJWTService(test.Secret, test.Method, mock)
			token, err := jwtService.GenerateToken(test.Input)

			if test.IsError {
				require.Error(t, err)
				require.True(t, errors.Is(err, err))
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Result, token)
			}
		})
	}
}

func TestJWTServiceVerifyToken(t *testing.T) {
	testCases := []TestCaseVerify{
		{
			Name:   "Success",
			Secret: "testkey",
			Method: jwt.SigningMethodHS256,
			Input:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoiVGVzdFVzZXIiLCJpZCI6InRlc3QifSwiZXhwIjoxNjg1OTEyNDAwLCJpYXQiOjE2ODUzMDc2MDB9.MknORRA4rFIqxbLk6C2JjQQQ_xIJ84rDQQMqV4l5Woo",
			Result: &model.Author{
				ID:       "test",
				Username: "TestUser",
			},
			IsError: false,
			Error:   nil,
		},
		{
			Name:   "Error Signature Invalid",
			Secret: "",
			Method: jwt.SigningMethodHS256,
			Input:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoiVGVzdFVzZXIiLCJpZCI6InRlc3QifSwiZXhwIjoxNjg1OTEyNDAwLCJpYXQiOjE2ODUzMDc2MDB9.YZtlV-onUbwMmPL0o5Jh6DPJ3OGb9o1UKeaULM7B3SA",
			Result: &model.Author{
				ID:       "test",
				Username: "TestUser",
			},
			IsError: true,
			Error:   jwt.ErrSignatureInvalid,
		},
		{
			Name:   "Error Invalid Sign Method",
			Secret: "testkey",
			Method: jwt.SigningMethodHS256,
			Input:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			Result: &model.Author{
				ID:       "test",
				Username: "TestUser",
			},
			IsError: true,
			Error:   model.ErrInvalidSignMethod,
		},
		{
			Name:   "Error Invalid Token",
			Secret: "testkey",
			Method: jwt.SigningMethodHS256,
			Input:  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoiVGVzdFVzZXIiLCJpZCI6InRlc3QifSwiZXhwIjoxNTAwMDAwMDAwLCJpYXQiOjE0MDAwMDAwMDB9.vyxAS4QNmxmPXFZmiI4zYya20L0PftIRKFnsCZWJ-B4",
			Result: &model.Author{
				ID:       "test",
				Username: "TestUser",
			},
			IsError: true,
			Error:   jwt.ErrTokenExpired,
		},
	}

	mock := new(FakeTimeController)
	mock.fixedTime = time.Date(2023, time.May, 29, 0, 0, 0, 0, time.UTC)

	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			jwtService := NewJWTService(test.Secret, test.Method, mock)
			author, err := jwtService.VerifyToken(test.Input)

			if test.IsError {
				require.Error(t, err)
				require.True(t, errors.Is(err, test.Error))
			} else {
				require.NoError(t, err)
				require.Equal(t, test.Result, author)
			}
		})
	}
}
