package entity

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Username string `json:"username"`
	Nama     string `json:"nama"`
	Role     string `json:"role"`
	Session  string `json:"session"`
	jwt.RegisteredClaims
}

type JwtCustomRefreshClaims struct {
	Username string `json:"username"`
	Nama     string `json:"nama"`
	Role     string `json:"role"`
	Session  string `json:"session"`
	jwt.RegisteredClaims
}

type GoogleClaims struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Athash        string `json:"at_hash"`
	Hd            string `json:"hd"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	Iat           int    `json:"iat"`
	Exp           int    `json:"exp"`
}
