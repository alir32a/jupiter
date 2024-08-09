package jwt

import "github.com/golang-jwt/jwt"

type Claim struct {
	jwt.StandardClaims
	Username string `json:"uname"`
}
