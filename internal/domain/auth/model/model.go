package model

import (
	"github.com/golang-jwt/jwt/v5"
)

// Claims is the JWT payload for an access token.
type Claims struct {
	UserId uint64 `json:"user_id"`
	jwt.RegisteredClaims
}
