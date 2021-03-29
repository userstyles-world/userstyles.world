package utils

import (
	"crypto/cipher"

	"golang.org/x/crypto/chacha20poly1305"
	"userstyles.world/config"
)

var AEAD cipher.AEAD

func InitalizeCrypto() {
	aead, err := chacha20poly1305.New([]byte(config.CRYPTO_KEY))
	if err != nil {
		panic("Cannot create AEAD chipher, due to " + err.Error())
	}
	AEAD = aead
}

func SealText(text string) []byte {
	nonce := RandStringBytesMaskImprSrcUnsafe(12)

	return AEAD.Seal(nil, []byte(nonce), []byte(text), nil)
}
