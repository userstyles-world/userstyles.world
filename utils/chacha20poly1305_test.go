package utils

import (
	"bytes"
	"testing"
)

func TestSimpleKey(t *testing.T) {
	InitalizeCrypto()

	jwt, err := NewJWTToken().
		SetClaim("email", "vednoc@usw.local").
		GetSignedString(VerifySigningKey)
	if err != nil {
		t.Error(err)
	}

	sealedText := SealText(jwt)
	unSealedText, err := OpenText(FastByteToString(sealedText))
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal([]byte(jwt), unSealedText) {
		t.Error("Originial and Unsealed aren't the same string.")
	}
}
