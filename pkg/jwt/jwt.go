package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
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
	UserID uint64
	Login  string
}

func CreateJWT(login string, userID uint64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenexp)),
		},
		Login:  login,
		UserID: userID,
	})
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func JWTPayload(tokenString string) (string, uint64, time.Time, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		return "", 0, time.Time{}, err
	}

	if !token.Valid {
		return "", 0, time.Time{}, ErrInvalidToken
	}

	return claims.Login, claims.UserID, claims.ExpiresAt.Time, nil
}
