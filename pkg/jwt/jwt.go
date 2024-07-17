package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

const (
	key      = "7OEdd8d8mOgLnIU9tLW5"
	tokenexp = time.Minute * 5
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Claims struct {
	jwt.RegisteredClaims
	Login string
}

func CreateJWT(login string) (string, error) {
	log.Printf("create token for login: %s", login)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenexp)),
		},
		Login: login,
	})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func JWTPayload(tokenString string) (string, time.Time, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return "", time.Time{}, err
	}

	if !token.Valid {
		return "", time.Time{}, ErrInvalidToken
	}
	log.Printf("parsed values login: %s", claims.Login)
	return claims.Login, claims.ExpiresAt.Time, nil
}
