package service

import (
	"analog-be/dto"
	"analog-be/entity"
	"analog-be/repository"
	"context"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	key    []byte
	parser *jwt.Parser
	repo   *repository.TokenRepository
}

func NewTokenService(repo *repository.TokenRepository) *TokenService {
	key, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_TOKEN"))

	if err != nil {
		panic(err)
	}

	return &TokenService{
		key: key,
		parser: jwt.NewParser(jwt.WithTimeFunc(func() time.Time {
			return time.Now().UTC()
		})),
		repo: repo,
	}
}

func (s *TokenService) sign(user *entity.User) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, dto.JwtClaims{
		Name:       user.Name,
		Generation: user.Generation,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "log.ana.st",
			ExpiresAt: jwt.NewNumericDate(user.JoinedAt.Add(time.Hour * 4)),
			IssuedAt:  jwt.NewNumericDate(user.JoinedAt),
			Subject:   strconv.FormatInt(user.ID, 10),
		},
	})

	signed, err := token.SignedString(s.key)

	if err != nil {
		return nil, err
	}

	return &signed, nil
}

func (s *TokenService) verify(tokenStr string) (*dto.JwtClaims, error) {
	token, err := s.parser.ParseWithClaims(tokenStr, &dto.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.key, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return token.Claims.(*dto.JwtClaims), nil
}

func (s *TokenService) generateRefreshToken(ctx context.Context, u *entity.User) (*entity.RefreshToken, error) {
	b := make([]byte, 64)
	_, err := rand.Read(b)

	if err != nil {
		return nil, err
	}

	bytes := sha512.Sum512(b)

	refreshToken := &entity.RefreshToken{
		Token:     hex.EncodeToString(bytes[:]),
		UserID:    u.ID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().AddDate(0, 0, 20),
	}

	err = s.repo.Create(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	return refreshToken, nil
}
