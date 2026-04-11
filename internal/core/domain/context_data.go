package domain

import "context"

const (
	UserIDKey string = "userID"
)

func GetUserID(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
