package helpers

import (
	"context"
	"crypto/rand"
	"encoding/base32"

	"golang.org/x/crypto/bcrypt"
)

func PasswordHash(passwd string) (string, error) {
	passwordBytes := []byte(passwd)
	hashedPasswordBytes, err := bcrypt.
		GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	return string(hashedPasswordBytes), SendTraceErrorToSentry(err)
}

func Generate2FASecretKey(ctx context.Context) (string, error) {
	randomBytes := make([]byte, 10)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes), nil
}
