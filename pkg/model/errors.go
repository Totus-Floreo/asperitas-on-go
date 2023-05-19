package model

import "errors"

// Немного не удобно использовать, по всей видимости из-за изначально неправильно выбора обработки ошибок

// type ErrorStack struct {
// 	MsgErrors []ErrorMessage `json:"errors"`
// }

// func NewErrorStack(location string, param string, value string, msg string) *ErrorStack {
// 	return &ErrorStack{
// 		MsgErrors: []ErrorMessage{
// 			ErrorMessage{
// 				Location: location,
// 				Param:    param,
// 				Value:    value,
// 				Msg:      msg,
// 			},
// 		},
// 	}
// }

// type ErrorMessage struct {
// 	Location string `json:"location"`
// 	Param    string `json:"param"`
// 	Value    string `json:"value"`
// 	Msg      string `json:"msg"`
// }

var (
	ErrUserNotFound       = errors.New("user doesn't exist")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrInvalidUrl         = errors.New(`{"errors":[{"location":"body","param":"url","value":"https:/broken.link/reallybad","msg":"is invalid"}]}`)
	ErrUserExist          = errors.New(`{"errors":[{"location":"body","param":"username","value":"testdata","msg":"already exists"}]}`)
	ErrPostNotFound       = errors.New(`{"message":"invalid post id"}`)
	ErrCommentNotFound    = errors.New(`{"message":"invalid comment id"}`)
	ErrInvalidToken       = errors.New("invalid token")
	ErrInvalidSignMethod  = errors.New("invalid sign method")
	ErrVoteNotFound       = errors.New("vote doesn't exist")
	ErrUnAuthorized       = errors.New(`{"message":"unuthorized"}`)
	ErrNullComment        = errors.New(`{"errors":[{"location":"body","param":"comment","msg":"is required"}]}`)
)
