package auth

type User struct {
	ID         string
	Email      string
	Name       string
	Nickname   string
	AvatarURL  string
	Provider   string
	ProviderID string
}
