package utils

import (
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"net/url"

	"github.com/form3tech-oss/jwt-go"
	"golang.org/x/crypto/chacha20poly1305"
	"userstyles.world/config"
)

var (
	AEAD                cipher.AEAD
	VerifySigningKey    = []byte(config.VERIFY_JWT_SIGNING_KEY)
	verifySigningMethod = "HS512"
)

func InitalizeCrypto() {
	aead, err := chacha20poly1305.NewX([]byte(config.CRYPTO_KEY))
	if err != nil {
		panic("Cannot create AEAD chipher, due to " + err.Error())
	}
	AEAD = aead
}

func SealText(text string) []byte {
	nonce := []byte(RandStringBytesMaskImprSrcUnsafe(AEAD.NonceSize()))

	return AEAD.Seal(nonce, nonce, []byte(text), nil)
}

func OpenText(encryptedMsg string) ([]byte, error) {
	// Split nonce and ciphertext.
	nonce, ciphertext := encryptedMsg[:AEAD.NonceSize()], encryptedMsg[AEAD.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	return AEAD.Open(nil, []byte(nonce), []byte(ciphertext), nil)

}

func verifyJwtKeyFunction(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != verifySigningMethod {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
	}
	return VerifySigningKey, nil
}

func PrepareText(text string) string {
	// We have to prepare the encrypted text for transport
	// Seal Text -> Base64 -> Path Escape

	sealedText := SealText(text)

	base64 := base64.StdEncoding.EncodeToString(sealedText)

	return url.PathEscape(base64)
}

func DecodePreparedText(preparedText string) (*jwt.Token, error) {
	// Now we have to reverse the process.
	// PathUnescape -> Decode base64 -> Unseal Text

	unescapedText, err := url.PathUnescape(preparedText)
	if err != nil {
		return nil, err
	}

	enryptedText, err := base64.StdEncoding.DecodeString(unescapedText)
	if err != nil {
		return nil, err
	}

	decryptedText, err := OpenText(FastByteToString(enryptedText))
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(FastByteToString(decryptedText), verifyJwtKeyFunction)
	if err != nil || !token.Valid {
		return nil, err
	}

	return token, nil

}
