package utils

import (
	"crypto/rand"
	"encoding/base64"
	"io"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateSessionId() string {
	id := make([]byte, 32) // 32 bytes = 256 bits
	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		panic("failed to generate session id")
	}
	return base64.RawURLEncoding.EncodeToString(id)
}
