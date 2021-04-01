package utils

import (
	"crypto/cipher"
	"encoding/base64"
	"fmt"

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
	nonce := RandStringBytesMaskImprSrcUnsafe(AEAD.NonceSize())

	return AEAD.Seal(nonce, nonce, S2b(text), nil)
}

func OpenText(encryptedMsg string) ([]byte, error) {
	// Split nonce and ciphertext.
	nonce, ciphertext := encryptedMsg[:AEAD.NonceSize()], encryptedMsg[AEAD.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	return AEAD.Open(nil, S2b(nonce), S2b(ciphertext), nil)

}

func verifyJwtKeyFunction(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != verifySigningMethod {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
	}
	return VerifySigningKey, nil
}

func PrepareText(text string) string {
	// We have to prepare the encrypted text for transport
	// Seal Text -> Base64(URL Version)
	sealedText := SealText(text)

	return EncodeToString(sealedText)
}

func DecodePreparedText(preparedText string) (*jwt.Token, error) {
	// Now we have to reverse the process.
	// Decode Base64(URL version) -> Unseal Text
	enryptedText, err := base64.URLEncoding.DecodeString(preparedText)
	if err != nil {
		return nil, err
	}

	decryptedText, err := OpenText(B2s(enryptedText))
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(B2s(decryptedText), verifyJwtKeyFunction)
	if err != nil || !token.Valid {
		return nil, err
	}

	return token, nil

}
