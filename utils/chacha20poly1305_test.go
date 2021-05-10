package utils

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/form3tech-oss/jwt-go"
)

func TestSimpleKey(t *testing.T) {
	InitalizeCrypto()

	jwtToken, err := NewJWTToken().
		SetClaim("email", "vednoc@usw.local").
		GetSignedString(VerifySigningKey)
	if err != nil {
		t.Error(err)
	}

	sealedText := SealText(jwtToken, AEAD_CRYPTO)
	unSealedText, err := OpenText(UnsafeString(sealedText), AEAD_CRYPTO)
	if err != nil {
		t.Error(err)
	}
	token, err := jwt.Parse(UnsafeString(unSealedText), VerifyJwtKeyFunction)
	if err != nil || !token.Valid {
		t.Error(err)
	}

	if !bytes.Equal(UnsafeBytes(jwtToken), unSealedText) {
		t.Error("Originial and Unsealed aren't the same string.")
	}
}

func benchamarkChaCha20Poly1305Seal(b *testing.B, buf []byte) {
	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SealText(UnsafeString(buf[:]), AEAD_CRYPTO)
	}
}

func benchamarkChaCha20Poly1305Open(b *testing.B, buf []byte) {
	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	ct := SealText(UnsafeString(buf[:]), AEAD_CRYPTO)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = OpenText(UnsafeString(ct[:]), AEAD_CRYPTO)
	}
}

func benchamarkPrepareText(b *testing.B, buf []byte) {
	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PrepareText(UnsafeString(buf[:]), AEAD_CRYPTO)
	}
}

func benchamarkDecodePreparedText(b *testing.B, buf []byte) {
	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	ct := PrepareText(UnsafeString(buf[:]), AEAD_CRYPTO)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecodePreparedText(ct, AEAD_CRYPTO)
	}
}

func BenchmarkPureChaCha20Poly1305(b *testing.B) {
	InitalizeCrypto()
	b.ResetTimer()
	for _, length := range []int{215, 1350, 8 * 1024} {
		b.Run("Open-"+strconv.Itoa(length)+"-X", func(b *testing.B) {
			benchamarkChaCha20Poly1305Open(b, make([]byte, length))
		})
		b.Run("Seal-"+strconv.Itoa(length)+"-X", func(b *testing.B) {
			benchamarkChaCha20Poly1305Seal(b, make([]byte, length))
		})
	}
}

func BenchmarkPrepareText(b *testing.B) {
	InitalizeCrypto()
	b.ResetTimer()
	for _, length := range []int{215, 1350, 8 * 1024} {
		b.Run("Prepare-"+strconv.Itoa(length), func(b *testing.B) {
			benchamarkPrepareText(b, make([]byte, length))
		})
		b.Run("Decode-"+strconv.Itoa(length), func(b *testing.B) {
			benchamarkDecodePreparedText(b, make([]byte, length))
		})
	}
}
