package util

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
)

// Check if we can create a JWT Token and sign it properly.
func TestGenerationToken(t *testing.T) {
	t.Parallel()
	token, err := NewJWT().
		SetClaim("sub", "test").
		SetExpiration(time.Now().Add(time.Minute * 5)).
		GetSignedString([]byte("secret"))
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("Token is empty")
	}
}

// Generate a Token and decrypt it again and check if the values are correct.
func TestDecryptionToken(t *testing.T) {
	t.Parallel()
	token, err := NewJWT().
		SetClaim("sub", "test").
		SetExpiration(time.Now().Add(time.Minute * 5)).
		GetSignedString([]byte("secret"))
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("Token is empty")
	}
	decrypted, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		t.Error(err)
	}
	// Check if claims types are []map[string]interface{}
	claims, ok := decrypted.Claims.(jwt.MapClaims)
	if !ok {
		t.Error("Claims are not a map")
	}

	if claims["sub"] != "test" {
		t.Error("Token decrypted with wrong values")
	}
}

// Test invalid expired token.
func TestExpiredToken(t *testing.T) {
	t.Parallel()
	token, err := NewJWT().
		SetClaim("sub", "test").
		SetExpiration(time.Now().Add(-time.Minute * 5)).
		GetSignedString([]byte("secret"))
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("Token is empty")
	}
	_, err = jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte("secret"), nil
	})
	if err == nil {
		t.Error("Expired token is valid")
	}
}

// Test invalid key signature.
func TestInvalidSignature(t *testing.T) {
	t.Parallel()
	token, err := NewJWT().
		SetClaim("sub", "test").
		SetExpiration(time.Now().Add(time.Minute * 5)).
		GetSignedString([]byte("secret"))
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("Token is empty")
	}
	_, err = jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte("secret2"), nil
	})
	if err == nil {
		t.Error("Invalid signature is valid")
	}
}

func TestInvalidSignature2(t *testing.T) {
	t.Parallel()
	token, err := NewJWT().
		SetClaim("sub", "test").
		SetExpiration(time.Now().Add(time.Minute * 5)).
		GetSignedString([]byte("secret"))
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Error("Token is empty")
	}
	_, err = jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return nil, errors.New("Some error")
	})
	if err == nil {
		t.Error("Invalid signature is valid")
	}
}

func TestMapclaimsVerifyIssuedAtInvalidTypeString(t *testing.T) {
	t.Parallel()

	mapClaims := jwt.MapClaims{
		"iat": "foo",
	}
	want := false
	got := mapClaims.VerifyIssuedAt(0, false)
	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}

func TestMapclaimsVerifyNotBeforeInvalidTypeString(t *testing.T) {
	t.Parallel()

	mapClaims := jwt.MapClaims{
		"nbf": "foo",
	}
	want := false
	got := mapClaims.VerifyNotBefore(0, false)
	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}

func TestMapclaimsVerifyExpiresAtInvalidTypeString(t *testing.T) {
	t.Parallel()

	mapClaims := jwt.MapClaims{
		"exp": "foo",
	}
	want := false
	got := mapClaims.VerifyExpiresAt(0, false)

	if want != got {
		t.Fatalf("Failed to verify claims, wanted: %v got %v", want, got)
	}
}
