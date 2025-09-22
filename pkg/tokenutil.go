package pkg

import (
	"e-klinik/internal/domain/entity"
	"e-klinik/utils"
	"errors"
	"fmt"
	"log"

	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(u entity.User, secret string, expiry int) (accessToken string, expire int64, err error) {
	now := time.Now()
	exp := now.Add(time.Minute * time.Duration(expiry))
	claims := &entity.JwtCustomClaims{
		Username: u.Username,
		Nama:     u.Nama,
		Role:     utils.DerefString(u.Role),
		Session:  u.Session,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "e-klink-track",
			Subject:   fmt.Sprint(u.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return t, now.Unix(), err
}

func CreateRefreshToken(u entity.User, secret string, expiry int) (refreshToken string, expire int64, err error) {
	now := time.Now()
	exp := now.Add(time.Hour * time.Duration(expiry))
	claimsRefresh := &entity.JwtCustomRefreshClaims{
		Username: u.Username,
		Nama:     u.Nama,
		Role:     utils.DerefString(u.Role),
		Session:  u.Session,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "e-klink-track",
			Subject:   fmt.Sprint(u.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claimsRefresh)
	rt, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}
	return rt, now.Unix(), err
}

func IsAuthorized(requestToken string, secret string) (bool, error) {
	token, err := jwt.ParseWithClaims(requestToken, &entity.JwtCustomRefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		// log detailed error
		log.Printf("JWT parse error: %v", err)

		// return a generalized error to the caller
		if errors.Is(err, jwt.ErrTokenExpired) {
			return false, fmt.Errorf("token is expired")
		}
		return false, fmt.Errorf("invalid token")
	}

	if !token.Valid {
		return false, fmt.Errorf("token is invalid")
	}

	return true, nil
}

func ExtractIDFromToken(requestToken string, secret string) (string, error) {
	claims := &entity.JwtCustomRefreshClaims{}

	token, err := jwt.ParseWithClaims(requestToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token is invalid")
	}

	return claims.RegisteredClaims.Subject, nil
}
func ExtractClaimsFromToken(requestToken string, secret string) (*entity.JwtCustomRefreshClaims, error) {
	claims := &entity.JwtCustomRefreshClaims{}

	token, err := jwt.ParseWithClaims(requestToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token is invalid")
	}

	// Return the claims object instead of just the subject.
	return claims, nil
}
