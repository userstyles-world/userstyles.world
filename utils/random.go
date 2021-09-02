package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomBytes returns securely generated random bytes.
// It's a helper function for crypto/rand.
func RandomBytes(n int) []byte {
	buffer := make([]byte, n)
	_, _ = rand.Read(buffer)
	return buffer
}

func RandomString(n int) string {
	src := RandomBytes(n)
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return UnsafeString(dst)
}
