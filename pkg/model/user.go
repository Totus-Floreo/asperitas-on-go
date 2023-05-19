package model

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
	GetUser(string) (*User, error)
	AddUser(*User) error
}
