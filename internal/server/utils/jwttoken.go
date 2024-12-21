package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

var (
	tokenTTL  = jwt.NewNumericDate(time.Now().Add(time.Hour * 3))
	secretKey = []byte("SuperSecretKey")
)

var ErrInvalidToken = errors.New("invalid token")

func CreateJWT(uid string) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: tokenTTL,
		},
		UserID: uid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

func ValidToken(tokentStr string) (string, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(tokentStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})
	if err != nil {
		return ``, err
	}
	if !token.Valid {
		return ``, ErrInvalidToken
	}
	return claims.UserID, nil
}
