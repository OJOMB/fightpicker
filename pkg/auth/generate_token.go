package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

type TokenGenerator struct {
}

// GenerateVerificationToken generates a secure random verification token.
func (tg *TokenGenerator) GenerateVerificationToken() (string, error) {
	b := make([]byte, 32) // 256 bits
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HashVerificationToken hashes the provided token using SHA-256 and returns the hash.
func (tg *TokenGenerator) HashVerificationToken(token string) ([]byte, error) {
	h := sha256.New()
	if _, err := h.Write([]byte(token)); err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
