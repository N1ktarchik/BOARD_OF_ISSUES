package domain

import "context"

type ContextKey string

const (
	UserIDKey ContextKey = "userID"
)

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
