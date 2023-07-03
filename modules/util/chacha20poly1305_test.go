package util

import (
	"bytes"
	"strconv"
	"testing"

	"github.com/golang-jwt/jwt"

	"userstyles.world/modules/config"
	"userstyles.world/utils"
)

func TestSimpleKey(t *testing.T) {
	t.Parallel()
	InitCrypto()

	jwtToken, err := utils.NewJWTToken().
		SetClaim("email", "vednoc@usw.local").
		GetSignedString(VerifySigningKey)
	if err != nil {
		t.Error(err)
	}

	scrambleConfig := &config.ScrambleSettings{
		StepSize:       3,
		BytesPerInsert: 2,
	}

	sealedText := sealText(jwtToken, AEADCrypto, scrambleConfig)
	unSealedText, err := openText(UnsafeString(sealedText), AEADCrypto, scrambleConfig)
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

func benchamarkChaCha20Poly1305Seal(b *testing.B, buf []byte, scrambleConfig *config.ScrambleSettings) {
	b.Helper()

	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sealText(UnsafeString(buf), AEADCrypto, scrambleConfig)
	}
}

func benchamarkChaCha20Poly1305Open(b *testing.B, buf []byte, scrambleConfig *config.ScrambleSettings) {
	b.Helper()

	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	ct := sealText(UnsafeString(buf), AEADCrypto, scrambleConfig)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = openText(UnsafeString(ct), AEADCrypto, scrambleConfig)
	}
}

func benchamarkPrepareText(b *testing.B, buf []byte, scrambleConfig *config.ScrambleSettings) {
	b.Helper()

	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = EncryptText(UnsafeString(buf), AEADCrypto, scrambleConfig)
	}
}

func benchamarkDecodePreparedText(b *testing.B, buf []byte, scrambleConfig *config.ScrambleSettings) {
	b.Helper()

	b.ReportAllocs()
	b.SetBytes(int64(len(buf)))

	ct := EncryptText(UnsafeString(buf), AEADCrypto, scrambleConfig)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = DecryptText(ct, AEADCrypto, scrambleConfig)
	}
}

func BenchmarkPureChaCha20Poly1305(b *testing.B) {
	InitCrypto()
	b.ResetTimer()

	scrambleConfig := &config.ScrambleSettings{
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
	InitCrypto()
	b.ResetTimer()

	scrambleConfig := &config.ScrambleSettings{
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

	nonce := "1124551523355"
	text := "ohnoweowfsdfsfdsfsd"

	dest := UnsafeString(scrambleNonce(UnsafeBytes(nonce), UnsafeBytes(text), 3, 3))
	if len(dest) == len(nonce)+len(text) && dest != "112ohn455owe152owf335sdf5sfdsfsd" {
		t.Error("Nonce scrambling failed.")
	}
}

func TestNonceDescrambling(t *testing.T) {
	t.Parallel()

	nonce := "1241312231412"
	text := "HellloBeautfikfuldasa"

	dest := scrambleNonce(UnsafeBytes(nonce), UnsafeBytes(text), 2, 3)

	if len(dest) == len(nonce)+len(text) && string(dest) != "124He131ll223lo141Be2autfikfuldasa" {
		t.Error("Nonce descrambling failed.")
	}

	// In production we know the Nonce of a specific hash, due to,
	// that AEAD is used. Which used a hard-coded length.
	descrambledNonce, descrambledText, err := descrambleNonce(dest, len(nonce), 2, 3)
	if err != nil {
		t.Error("Couldn't descramble, errored:", err)
	}

	if string(descrambledNonce) != nonce {
		t.Error("Couldn't descramble nonce")
	}

	if string(descrambledText) != text {
		t.Error("Couldn't descramble text")
	}
}

func TestNonceDescramblingInsaneConfig(t *testing.T) {
	t.Parallel()

	nonce := "123132131312312312312312312312313123123123123123"
	text := "HellloBeautfikfuldasa"

	dest := scrambleNonce(UnsafeBytes(nonce), UnsafeBytes(text), 9, 20)

	if len(dest) == len(nonce)+len(text) &&
		string(dest) != "12313213131231231231HellloBea23123123123131231231utfikfuld23123123asa" {
		t.Error("Nonce descrambling failed.")
	}

	// In production we know the Nonce of a specific hash, due to,
	// that AEAD is used. Which used a hard-coded length.
	descrambledNonce, descrambledText, err := descrambleNonce(dest, len(nonce), 9, 20)
	if err != nil {
		t.Error("Couldn't descramble, errored:", err)
	}

	if string(descrambledNonce) != nonce {
		t.Error("Couldn't descramble nonce")
	}

	if string(descrambledText) != text {
		t.Error("Couldn't descramble text")
	}
}

func TestNonceDescramblingPerfectStop(t *testing.T) {
	t.Parallel()

	nonce := "1234567890"
	text := "abcdefghijklmnopqr"

	dest := scrambleNonce(UnsafeBytes(nonce), UnsafeBytes(text), 3, 1)

	if len(dest) == len(nonce)+len(text) && string(dest) != "1abc2def3ghi4jkl5mno6pqr7890" {
		t.Error("Nonce descrambling failed.")
	}

	// In production we know the Nonce of a specific hash, due to,
	// that AEAD is used. Which used a hard-coded length.
	descrambledNonce, descrambledText, err := descrambleNonce(dest, len(nonce), 3, 1)
	if err != nil {
		t.Error("Couldn't descramble, errored:", err)
	}

	if string(descrambledNonce) != nonce {
		t.Error("Couldn't descramble nonce")
	}

	if string(descrambledText) != text {
		t.Error("Couldn't descramble text")
	}
}

func TestNonceDescramblingWithOverflow(t *testing.T) {
	t.Parallel()

	nonce := "124131223141274483127131231"
	text := "HellloBeautfikfuldasa"

	dest := scrambleNonce(UnsafeBytes(nonce), UnsafeBytes(text), 2, 1)

	if string(dest) != "1He2ll4lo1Be3au1tf2ik2fu3ld1as4a1274483127131231" {
		t.Error("Nonce descrambling failed.")
	}

	// In production we know the Nonce of a specific hash, due to,
	// that AEAD is used. Which used a hard-coded length.
	descrambledNonce, descrambledText, err := descrambleNonce(dest, len(nonce), 2, 1)
	if err != nil {
		t.Error("Couldn't descramble, errored:", err)
	}

	if string(descrambledNonce) != nonce {
		t.Error("Couldn't descramble nonce")
	}

	if string(descrambledText) != text {
		t.Error("Couldn't descramble text")
	}
}

func TestNonceDescramblingOnIncorrectInput(t *testing.T) {
	t.Parallel()

	dest := []byte("helloI'mMaliciousInput")
	_, _, err := descrambleNonce(dest, 24, 4, 1)

	if err == nil {
		t.Error("Descrambling should fail on incorrect input")
	}

	dest = []byte("hello")
	_, _, err = descrambleNonce(dest, 24, 4, 1)

	if err == nil {
		t.Error("Descrambling should fail on incorrect input")
	}

	dest = []byte("22")
	_, _, err = descrambleNonce(dest, 2, 4, 1)

	if err == nil {
		t.Error("Descrambling should fail on incorrect input")
	}

	dest = []byte("333")
	_, _, err = descrambleNonce(dest, 2, 4, 1)
	if err != nil {
		t.Error("Descrambling should't fall")
	}

	dest = []byte("4444")
	_, _, err = descrambleNonce(dest, 3, 1, 1)
	if err != nil {
		t.Error("Descrambling should't fall")
	}

	dest = []byte("55555")
	_, _, err = descrambleNonce(dest, 4, 3, 4)
	if err != nil {
		t.Error("Descrambling should't fall")
	}
}
