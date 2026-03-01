package dto

import "github.com/golang-jwt/jwt/v5"

type JwtClaims struct {
	Name       string
	Generation uint16
	jwt.RegisteredClaims
}
