package auth

import (
	"N1ktarchik/Board_of_issues/internal/core/domain"
	"N1ktarchik/Board_of_issues/internal/core/errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *JWTService) CreateJWT(userID string) (string, error) {
	secret := s.secretKey

	claims := domain.Claims{
		ID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.liveTimeMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "default-issuer",
		},
	}

	JWT, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		s.log.Error("failed to sign jwt token", slog.String("user_id", userID), slog.Any("err", err))
		return "", err
	}

	s.log.Debug("jwt token created successfully", slog.String("user_id", userID))

	return JWT, nil
}

func (s *JWTService) ValidateJWT(JWT string) (*domain.Claims, error) {

	claims := &domain.Claims{}

	token, err := jwt.ParseWithClaims(JWT, claims, func(t *jwt.Token) (interface{}, error) {

		if t.Method != jwt.SigningMethodHS256 {
			s.log.Warn("invalid jwt signing method", slog.Any("method", t.Header["alg"]))
			return nil, errors.JWTTokenNotValid()
		}

		return s.secretKey, nil
	})

	if err != nil {
		s.log.Warn("jwt validation failed", slog.Any("err", err))
		return nil, err
	}

	if token.Valid {
		return claims, nil
	}

	s.log.Warn("jwt token is invalid")
	return nil, errors.JWTTokenNotValid()
}

func (s *JWTService) GetUserIdFromJWT(JWT string) (string, error) {
	claims, err := s.ValidateJWT(JWT)
	if err != nil {
		return "", err
	}

	return claims.ID, nil
}
