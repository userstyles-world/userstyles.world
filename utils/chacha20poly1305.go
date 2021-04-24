package utils

import (
	"crypto/cipher"
	"errors"
	"fmt"

	"github.com/form3tech-oss/jwt-go"
	"golang.org/x/crypto/chacha20poly1305"
	"userstyles.world/config"
)

var (
	AEAD_CRYPTO      cipher.AEAD
	AEAD_OAUTH       cipher.AEAD
	AEAD_OAUTHP      cipher.AEAD
	VerifySigningKey = []byte(config.VERIFY_JWT_SIGNING_KEY)
	OAuthPSigningKey = []byte(config.OAUTHP_JWT_SIGNING_KEY)
	signingMethod    = "HS512"
)

func InitalizeCrypto() {
	var aead cipher.AEAD
	var err error

	aead, err = chacha20poly1305.NewX([]byte(config.CRYPTO_KEY))
	if err != nil {
		panic("Cannot create AEAD_CRYPTO chipher, due to " + err.Error())
	}
	AEAD_CRYPTO = aead

	aead, err = chacha20poly1305.NewX([]byte(config.OAUTH_KEY))
	if err != nil {
		panic("Cannot create AEAD_OAUTH chipher, due to " + err.Error())
	}
	AEAD_OAUTH = aead

	aead, err = chacha20poly1305.NewX([]byte(config.OAUTH_KEY))
	if err != nil {
		panic("Cannot create AEAD_OAUTHP chipher, due to " + err.Error())
	}
	AEAD_OAUTHP = aead
}

func SealText(text string, aead cipher.AEAD) []byte {
	nonce := RandStringBytesMaskImprSrcUnsafe(aead.NonceSize())

	return aead.Seal(nonce, nonce, S2b(text), nil)
}

func OpenText(encryptedMsg string, aead cipher.AEAD) ([]byte, error) {
	if len(encryptedMsg) < aead.NonceSize() {
		return nil, errors.New("message too small")
	}
	// Split nonce and ciphertext.
	nonce, ciphertext := encryptedMsg[:aead.NonceSize()], encryptedMsg[aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	return aead.Open(nil, S2b(nonce), S2b(ciphertext), nil)

}

func VerifyJwtKeyFunction(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != signingMethod {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
	}
	return VerifySigningKey, nil
}

func OAuthPJwtKeyFunction(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != signingMethod {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
	}
	return OAuthPSigningKey, nil
}

func PrepareText(text string, aead cipher.AEAD) string {
	// We have to prepare the encrypted text for transport
	// Seal Text -> Base64(URL Version)
	sealedText := SealText(text, aead)

	return EncodeToString(sealedText)
}

func DecodePreparedText(preparedText string, aead cipher.AEAD) (string, error) {
	// Now we have to reverse the process.
	// Decode Base64(URL version) -> Unseal Text
	enryptedText, err := DecodeString(preparedText)
	if err != nil {
		return "", err
	}

	decryptedText, err := OpenText(B2s(enryptedText), aead)
	if err != nil {
		return "", err
	}

	return B2s(decryptedText), nil

}
