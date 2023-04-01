package authorizing

import "github.com/markbates/goth"

type UserView struct {
	ID        string
	Name      string
	Nickname  string
	Email     string
	AvatarURL string
	Provider  string
}

func NewUserView(g goth.User) UserView {
	return UserView{
		ID:        g.UserID,
		Name:      g.Name,
		Nickname:  g.NickName,
		Email:     g.Email,
		AvatarURL: g.AvatarURL,
		Provider:  g.Provider,
	}
}
