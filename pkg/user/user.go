package user

import "github.com/markbates/goth"

type User struct {
	ID    string
	Name  string
	Email string
}

func FromGothUser(g goth.User) User {
	usr := User{
		Name:  "",
		Email: "",
	}
	return usr
}
