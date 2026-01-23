package pkg

import (
	"analog-be/entity"
	"context"
)

type contextKey string

const (
	UserIDKey       contextKey = "userID"
	SessionTokenKey contextKey = "sessionToken"
)

func GetUserID(ctx context.Context) (entity.ID, bool) {
	userID, ok := ctx.Value(UserIDKey).(entity.ID)
	return userID, ok
}

func GetSessionToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(SessionTokenKey).(string)
	return token, ok
}
