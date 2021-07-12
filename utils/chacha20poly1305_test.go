package utils

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/form3tech-oss/jwt-go"
	"userstyles.world/modules/config"
)

func TestSimpleKey(t *testing.T) {
	t.Parallel()

	InitalizeCrypto()

	jwtToken, err := NewJWTToken().
		SetClaim("email", "vednoc@usw.local").
		GetSignedString(VerifySigningKey)
	if err != nil {
		t.Error(err)
	}

	scrambleConfig := &config.NonceScramblingConfig{
		StepSize:       3,
		BytesPerInsert: 2,
	}

	sealedText := SealText(jwtToken, AEAD_CRYPTO, scrambleConfig)
	unSealedText, err := OpenText(UnsafeString(sealedText), AEAD_CRYPTO, scrambleConfig)
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

func benchamarkChaCha20Poly1305Seal(b *testing.B, buf []byte, scrambleConfig *config.NonceScramblingConfig) {
	b.Helper()

	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = SealText(UnsafeString(buf), AEAD_CRYPTO, scrambleConfig)
	}
}

func benchamarkChaCha20Poly1305Open(b *testing.B, buf []byte, scrambleConfig *config.NonceScramblingConfig) {
	b.Helper()

	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	ct := SealText(UnsafeString(buf), AEAD_CRYPTO, scrambleConfig)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = OpenText(UnsafeString(ct), AEAD_CRYPTO, scrambleConfig)
	}
}

func benchamarkPrepareText(b *testing.B, buf []byte, scrambleConfig *config.NonceScramblingConfig) {
	b.Helper()

	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = PrepareText(UnsafeString(buf), AEAD_CRYPTO, scrambleConfig)
	}
}

func benchamarkDecodePreparedText(b *testing.B, buf []byte, scrambleConfig *config.NonceScramblingConfig) {
	b.Helper()

	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	ct := PrepareText(UnsafeString(buf), AEAD_CRYPTO, scrambleConfig)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecodePreparedText(ct, AEAD_CRYPTO, scrambleConfig)
	}
}

func BenchmarkPureChaCha20Poly1305(b *testing.B) {
	InitalizeCrypto()
	b.ResetTimer()

	scrambleConfig := &config.NonceScramblingConfig{
		StepSize:       2,
		BytesPerInsert: 4,
	}

	for _, length := range []int{215, 1350, 8 * 1024} {
		b.Run("Open-"+strconv.Itoa(length)+"-X", func(b *testing.B) {
			benchamarkChaCha20Poly1305Open(b, make([]byte, length), scrambleConfig)
		})
		b.Run("Seal-"+strconv.Itoa(length)+"-X", func(b *testing.B) {
			benchamarkChaCha20Poly1305Seal(b, make([]byte, length), scrambleConfig)
		})
	}
}

func BenchmarkPrepareText(b *testing.B) {
	InitalizeCrypto()
	b.ResetTimer()

	scrambleConfig := &config.NonceScramblingConfig{
		StepSize:       2,
		BytesPerInsert: 4,
	}

	for _, length := range []int{215, 1350, 8 * 1024} {
		b.Run("Prepare-"+strconv.Itoa(length), func(b *testing.B) {
			benchamarkPrepareText(b, make([]byte, length), scrambleConfig)
		})
		b.Run("Decode-"+strconv.Itoa(length), func(b *testing.B) {
			benchamarkDecodePreparedText(b, make([]byte, length), scrambleConfig)
		})
	}
}

func TestNonceEncoding(t *testing.T) {
	t.Parallel()
	InitalizeCrypto()

	nonce := "1124551523355"
	text := "ohnoweowfsdfsfdsfsd"

	dest := UnsafeString(ScrambleNonce(UnsafeBytes(nonce), UnsafeBytes(text), 3, 3))
	if len(dest) == len(nonce)+len(text) && dest != "112ohn455owe152owf335sdf5sfdsfsd" {
		t.Error("Nonce scrambling failed.")
	}
}

func TestNonceDescrambling(t *testing.T) {
	t.Parallel()
	InitalizeCrypto()

	nonce := "1241312231412"
	text := "HellloBeautfikfuldasa"

	dest := ScrambleNonce(UnsafeBytes(nonce), UnsafeBytes(text), 2, 3)

	if len(dest) == len(nonce)+len(text) && string(dest) != "124He131ll223lo141Be2autfikfuldasa" {
		t.Error("Nonce descrambling failed.")
	}

	// In production we know the Nonce of a specific hash, due to,
	// that AEAD is used. Which used a hard-coded length.
	descrambledNonce, descrambledText := DescrambleNonce(dest, len(nonce), 2, 3)

	if string(descrambledNonce) != nonce {
		t.Error("Couldn't descramble nonce")
	}

	if string(descrambledText) != text {
		t.Error("Couldn't descramble text")
	}
}

func TestNonceDescramblingWithOverflow(t *testing.T) {
	t.Parallel()
	InitalizeCrypto()

	nonce := "124131223141274483127131231"
	text := "HellloBeautfikfuldasa"

	dest := ScrambleNonce(UnsafeBytes(nonce), UnsafeBytes(text), 2, 1)

	if string(dest) != "1He2ll4lo1Be3au1tf2ik2fu3ld1as4a1274483127131231" {
		t.Error("Nonce descrambling failed.")
	}

	// In production we know the Nonce of a specific hash, due to,
	// that AEAD is used. Which used a hard-coded length.
	descrambledNonce, descrambledText := DescrambleNonce(dest, len(nonce), 2, 1)

	if string(descrambledNonce) != nonce {
		t.Error("Couldn't descramble nonce")
	}

	if string(descrambledText) != text {
		t.Error("Couldn't descramble text")
	}
}
