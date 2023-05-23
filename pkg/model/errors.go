package model

import (
	"errors"

	"encoding/json"
)

var (
	ErrPostNotFoundHTTP = errors.New(`{"message": "post not found"}`)

	ErrPostInvalidHTTP         = errors.New(`{"message":"invalid post id"}`)
	ErrPostCategoryInvalidHTTP = errors.New(`{"message":"invalid post category"}`)
	ErrCommentInvalidHTTP      = errors.New(`{"message":"invalid comment id"}`)
	ErrUserInvalidHTTP         = errors.New(`{"message":"invalid user name"}`)
	ErrInvalidCredentialsHTTP  = errors.New(`{"message": "invalid username or password"}`)

	// HTTPErrNullComment = errors.New(`{"errors":[{"location":"body","param":"comment","msg":"is required"}]}`)

	// HTTPErrUserExist = errors.New(`{"errors":[{"location":"body","param":"username","value":"ads","msg":"already exists"}]}`)

	// HTTPErrInvalidUrl = errors.New(`{"errors":[{"location":"body","param":"url","value":"https:/broken.link/reallybad","msg":"is invalid"}]}`)

	ErrUnAuthorizedHTTP = errors.New(`{"message":"unuthorized"}`)
)

type ErrorStack struct {
	MsgErrors []ErrorMessage `json:"errors"`
}

func NewErrorStack(location string, param string, value string, msg string) (string, error) {
	msgErr := ErrorStack{
		MsgErrors: []ErrorMessage{
			{
				Location: location,
				Param:    param,
				Value:    value,
				Msg:      msg,
			},
		},
	}
	httpErr, err := json.Marshal(msgErr)
	if err != nil {
		return "", err
	}
	return string(httpErr), nil
}

type ErrorMessage struct {
	Location string `json:"location"`
	Param    string `json:"param"`
	Value    string `json:"value,omitempty"`
	Msg      string `json:"msg"`
}

var (
	ErrUserNotFound    = errors.New("user doesn't exist")
	ErrPostNotFound    = errors.New("post doesn't exist")
	ErrCommentNotFound = errors.New("comment doesn't exist")
	ErrVoteNotFound    = errors.New("vote doesn't exist")

	ErrUserExist = errors.New("user already exist")

	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidUrl         = errors.New("invalid url")
	ErrInvalidSignMethod  = errors.New("invalid sign method")

	ErrUnAuthorized = errors.New("unuthorized")
)
