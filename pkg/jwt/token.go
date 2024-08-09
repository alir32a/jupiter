package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"time"
)

func CreateToken(username string, secret string, expireTime time.Duration) (string, error) {
	now := time.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add(expireTime).Unix(),
			IssuedAt:  now.Unix(),
		},
		Username: username,
	})

	return token.SignedString([]byte(secret))
}

func ParseToken(token, secret string) (*Claim, error) {
	t, err := jwt.ParseWithClaims(token, &Claim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if t.Valid {
		return t.Claims.(*Claim), nil
	}

	return nil, errors.New("token is invalid")
}
