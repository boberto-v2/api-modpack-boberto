package models

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}
