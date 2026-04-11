package domain

import (
	"crypto/sha256"
	"encoding/hex"

	"Board_of_issuses/internal/core/errors"
)

const (
	MinPasswordLength int = 6
	MaxPasswordLength int = 30
)

func Hash(password string) (string, error) {

	if len(password) <= MinPasswordLength {
		return "", errors.TooShortPassword()
	}

	if len(password) > MaxPasswordLength {
		return "", errors.TooLongPassword()
	}

	hashPassword := sha256.Sum256([]byte(password))

	return hex.EncodeToString(hashPassword[:]), nil
}

func Compare(password, trueHashPassword string) bool {
	hashPassword, err := Hash(password)
	if err != nil {
		return false
	}

	return hashPassword == trueHashPassword
}
