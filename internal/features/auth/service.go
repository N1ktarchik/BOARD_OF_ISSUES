package auth

import (
	"log/slog"
)

type JWTService struct {
	secretKey       []byte
	liveTimeMinutes int
	log             *slog.Logger
}

func CreateJWTService(secret string, log *slog.Logger, liveTimeMinutes int) *JWTService {
	return &JWTService{
		secretKey:       []byte(secret),
		log:             log,
		liveTimeMinutes: liveTimeMinutes,
	}
}
