package pkg

import (
	"e-klinik/internal/domain/entity"
	"errors"
	"fmt"

	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

func CreateAccessToken(u entity.User, secret string, expiryMinutes int) (string, int64, error) {
	now := time.Now().UTC()
	exp := now.Add(time.Minute * time.Duration(expiryMinutes))

	claims := &entity.JwtCustomRefreshClaims{
		Username: u.Username,
		Nama:     u.Nama,
		Role:     u.Role,
		Session:  u.Session,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "e-klink-track",
			Subject:   fmt.Sprint(u.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}

	return signed, exp.Unix(), nil
}

func CreateRefreshToken(u entity.User, secret string, expiryHours int) (string, int64, error) {
	now := time.Now().UTC()
	exp := now.Add(time.Hour * time.Duration(expiryHours))

	claims := &entity.JwtCustomRefreshClaims{
		Username: u.Username,
		Nama:     u.Nama,
		Role:     u.Role,
		Session:  u.Session,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "e-klink-track",
			Subject:   fmt.Sprint(u.ID),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", 0, err
	}

	return signed, exp.Unix(), nil
}

func IsAuthorized(tokenStr, secret string) (bool, *entity.JwtCustomRefreshClaims, error) {
	claims := &entity.JwtCustomRefreshClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		// JWT v5: cek token expired
		if errors.Is(err, jwt.ErrTokenExpired) {
			return false, claims, fmt.Errorf("token expired at %v", claims.ExpiresAt.Time)
		}
		return false, nil, fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return false, claims, fmt.Errorf("token invalid")
	}

	return true, claims, nil
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
