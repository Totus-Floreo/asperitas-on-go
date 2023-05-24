package model

import "context"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Author struct {
	Username string `json:"username"`
	ID       string `json:"id"`
}

type IUserStorage interface {
	GetUser(context.Context, string) (*User, error)
	AddUser(context.Context, *User) error
}
