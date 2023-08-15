package models

import "github.com/rs/xid"

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewUser(username string, password string) *User {
	id := xid.New().String()
	return &User{
		ID:       id,
		Username: username,
		Password: password,
	}
}
