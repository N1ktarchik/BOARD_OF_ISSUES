package server

import (
	"Board_of_issuses/internal/core/domain"
	"Board_of_issuses/internal/core/errors"
	resp "Board_of_issuses/internal/core/transport/response"
	"context"
	"log/slog"
	"net/http"
	"strings"
)

const (
	UserIDKey string = "userID"
)

type AuthService interface {
	CreateJWT(userID string) (string, error)
	GetUserIdFromJWT(JWT string) (string, error)
	ValidateJWT(JWT string) (*domain.Claims, error)
}

type MiddleWare struct {
	authService AuthService
	log         *slog.Logger
}

func NewMiddleWare(authService AuthService, log *slog.Logger) *MiddleWare {
	return &MiddleWare{
		authService: authService,
		log:         log,
	}
}

func (m *MiddleWare) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.log.Warn("auth middleware: missing authorization header")
			resp.RespondWithError(w, errors.BadRequest())
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			m.log.Warn("auth middleware: invalid authorization format", slog.String("header", authHeader))
			resp.RespondWithError(w, errors.BadRequest())
			return
		}

		tokenString := parts[1]

		claims, err := m.authService.ValidateJWT(tokenString)
		if err != nil {
			m.log.Warn("auth middleware: jwt validation failed", slog.Any("err", err))
			resp.RespondWithError(w, err)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.ID)

		m.log.Debug("auth middleware: user authenticated", slog.String("user_id", claims.ID))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
