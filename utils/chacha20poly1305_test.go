package utils_test

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/form3tech-oss/jwt-go"
	"userstyles.world/utils"
)

func TestSimpleKey(t *testing.T) {
	utils.InitalizeCrypto()

	jwtToken, err := utils.NewJWTToken().
		SetClaim("email", "vednoc@usw.local").
		GetSignedString(utils.VerifySigningKey)
	if err != nil {
		t.Error(err)
	}

	sealedText := utils.SealText(jwtToken, utils.AEAD_CRYPTO)
	unSealedText, err := utils.OpenText(utils.UnsafeString(sealedText), utils.AEAD_CRYPTO)
	if err != nil {
		t.Error(err)
	}
	token, err := jwt.Parse(utils.UnsafeString(unSealedText), utils.VerifyJwtKeyFunction)
	if err != nil || !token.Valid {
		t.Error(err)
	}

	if !bytes.Equal(utils.UnsafeBytes(jwtToken), unSealedText) {
		t.Error("Originial and Unsealed aren't the same string.")
	}
}

func benchamarkChaCha20Poly1305Seal(b *testing.B, buf []byte) {
	b.Helper()

	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.SealText(utils.UnsafeString(buf[:]), utils.AEAD_CRYPTO)
	}
}

func benchamarkChaCha20Poly1305Open(b *testing.B, buf []byte) {
	b.Helper()

	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	ct := utils.SealText(utils.UnsafeString(buf[:]), utils.AEAD_CRYPTO)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = utils.OpenText(utils.UnsafeString(ct[:]), utils.AEAD_CRYPTO)
	}
}

func benchamarkPrepareText(b *testing.B, buf []byte) {
	b.Helper()

	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = utils.PrepareText(utils.UnsafeString(buf[:]), utils.AEAD_CRYPTO)
	}
}

func benchamarkDecodePreparedText(b *testing.B, buf []byte) {
	b.Helper()

	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	ct := utils.PrepareText(utils.UnsafeString(buf[:]), utils.AEAD_CRYPTO)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = utils.DecodePreparedText(ct, utils.AEAD_CRYPTO)
	}
}

func BenchmarkPureChaCha20Poly1305(b *testing.B) {
	utils.InitalizeCrypto()
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
	utils.InitalizeCrypto()
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
