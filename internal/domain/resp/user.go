package resp

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Nama         string `json:"nama"`
	Role         string `json:"role"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Exp          int64  `json:"expire"`
}

type RefreshToken struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp"`
}
