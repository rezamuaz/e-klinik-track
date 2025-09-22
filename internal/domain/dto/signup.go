package dto

type SignupRequest struct {
	Name      string   `json:"name" binding:"required"`
	Email     string   `json:"email" binding:"required,email"`
	Password  string   `json:"password" binding:"required"`
	Role      []string `json:"role" binding:"required"`
	Place     string   `json:"place"`
	Room      string   `json:"room"`
	PlaceName string   `json:"place_name"`
	RoomName  string   `json:"room_name"`
}

type SignupResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
