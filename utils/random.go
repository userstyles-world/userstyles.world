package utils

import (
	"crypto/rand"
)

// GenerateRandomBytes returns securely generated random bytes.
// It's a helper function for crypto/rand.
func RandomString(n int) []byte {
	buffer := make([]byte, n)
	_, _ = rand.Read(buffer)
	return buffer
}
