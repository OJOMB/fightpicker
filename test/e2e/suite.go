package e2e

import (
	"math/rand"

	"github.com/oapi-codegen/runtime/types"
)

const (
	testDomain = "http://localhost:8080"

	createUserURL = "/api/v1/users"
)

func ptrString(s string) *string {
	return &s
}

func ptrOAPIEmail(email types.Email) *types.Email {
	return &email
}

func newRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func newRandomEmail() string {
	return newRandomString(10) + "@example.com"
}

func newRandomUsername() string {
	return newRandomString(8)
}
