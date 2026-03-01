package dto

import (
	"analog-be/entity"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaims struct {
	Name       string
	Generation uint16
	jwt.RegisteredClaims
}

type RefreshTokenResponse struct {
	Token        *string
	RefreshToken *entity.RefreshToken
}
