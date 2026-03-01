package pkg

import (
	"analog-be/dto"
	"analog-be/entity"
	"context"
)

type contextKey string

const (
	UserClaims contextKey = "userClaims"
	UserID     contextKey = "userID"
)

func GetUserClaims(ctx context.Context) (dto.JwtClaims, bool) {
	userClaims, ok := ctx.Value(UserClaims).(dto.JwtClaims)
	return userClaims, ok
}

func GetUserID(ctx context.Context) (entity.ID, bool) {
	userID, ok := ctx.Value(UserClaims).(entity.ID)
	return userID, ok
}
