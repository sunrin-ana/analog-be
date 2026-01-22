package pkg

import "context"

type contextKey string

const (
	UserIDKey       contextKey = "userID"
	SessionTokenKey contextKey = "sessionToken"
)

func GetUserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	return userID, ok
}

func GetSessionToken(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(SessionTokenKey).(string)
	return token, ok
}
