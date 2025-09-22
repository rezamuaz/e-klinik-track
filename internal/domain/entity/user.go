package entity

const (
	CollectionUser = "users"
)

type User struct {
	ID       string  `json:"id"`
	Username string  `json:"username"`
	Nama     string  `json:"nama"`
	Role     *string `json:"role"`
	Session  string  `json:"session,omitempty"`
}
type UserCreated struct {
	AppID        string   `json:"app_id"`
	Email        string   `json:"email"`
	Picture      string   `json:"picture"`
	Role         []string `json:"role"`
	GivenName    string   `json:"given_name"`
	FamilyName   string   `json:"family_name"`
	Name         string   `json:"name"`
	RefreshToken string   `json:"refresh_token"`
}
