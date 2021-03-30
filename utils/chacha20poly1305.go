package utils

import (
	"crypto/cipher"
	"fmt"
	"math/rand"

	"github.com/form3tech-oss/jwt-go"
	"golang.org/x/crypto/chacha20poly1305"
	"userstyles.world/config"
)

var (
	AEAD                cipher.AEAD
	VerifySigningKey    = []byte(config.VERIFY_JWT_SIGNING_KEY)
	VerifySigningMethod = "HS512"
)

func InitalizeCrypto() {
	aead, err := chacha20poly1305.NewX([]byte(config.CRYPTO_KEY))
	if err != nil {
		panic("Cannot create AEAD chipher, due to " + err.Error())
	}
	AEAD = aead
}

func SealText(text string) []byte {
	nonce := make([]byte, AEAD.NonceSize(), AEAD.NonceSize()+len(text)+AEAD.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		panic(err)
	}

	return AEAD.Seal(nonce, nonce, []byte(text), nil)
}

func OpenText(encryptedMsg string) ([]byte, error) {
	// Split nonce and ciphertext.
	nonce, ciphertext := encryptedMsg[:AEAD.NonceSize()], encryptedMsg[AEAD.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	return AEAD.Open(nil, []byte(nonce), []byte(ciphertext), nil)

}

func VerifyJwtKeyFunction(t *jwt.Token) (interface{}, error) {
	if t.Method.Alg() != VerifySigningMethod {
		return nil, fmt.Errorf("unexpected jwt signing method=%v", t.Header["alg"])
	}
	return VerifySigningKey, nil
}
