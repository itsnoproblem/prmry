package auth

type User struct {
	ID         string
	Name       string
	Nickname   string
	Email      string
	AvatarURL  string
	Provider   string
	ProviderID string
}
