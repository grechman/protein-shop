package api

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, 100, 64*1024, 4, 32)
	combined := append(salt, hash...)
	return base64.StdEncoding.EncodeToString(combined), nil
}

func VerifyPassword(hashed, password string) (bool, error) {
	decoded, err := base64.StdEncoding.DecodeString(hashed)
	if err != nil {
		return false, err
	}
	salt := decoded[:16]
	hash := decoded[16:]

	newHash := argon2.IDKey([]byte(password), salt, 100, 64*1024, 4, 32)
	return string(hash) == string(newHash), nil
}
