package util

import (
	"crypto/rand"

	"github.com/valyala/bytebufferpool"
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
	return encodeToHex(src)
}

const hextable = "0123456789abcdef"

func encodeToHex(src []byte) string {
	dst := bytebufferpool.Get()
	defer bytebufferpool.Put(dst)
	for _, v := range src {
		_ = dst.WriteByte(hextable[v>>4])
		_ = dst.WriteByte(hextable[v&0x0f])
	}
	return UnsafeString(dst.B)
}
