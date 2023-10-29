package utils

import (
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("Error while hashing password", err)
		return "", err
	}

	return string(bytes), nil
}

func PasswordMatches(password string, dbHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(dbHash), []byte(password))
	if err != nil {
		slog.Error("Error while comparing hash and password.", err)
	}

	return err == nil
}
