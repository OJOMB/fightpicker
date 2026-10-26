package e2e

import "math/rand"

func NewRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func NewRandomEmail() string {
	return NewRandomString(10) + "@example.com"
}

func NewRandomUsername() string {
	return NewRandomString(8)
}
