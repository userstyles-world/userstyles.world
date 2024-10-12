package util

import (
	"strings"
	"testing"
)

const pw = "UserStyles.world"

func TestHashPassword(t *testing.T) {
	t.Parallel()

	got, err := HashPassword(pw, &c.Secrets)
	if err != nil {
		t.Fatalf("bcrypt failed to hash a password: %s", err)
	}

	exp := "$2a$10"
	if !strings.HasPrefix(got, exp) {
		t.Errorf("%q should have suffix %q", got, exp)
	}
}

func BenchmarkHashPassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = HashPassword(pw, &c.Secrets)
	}
}

func TestVerifyPassword(t *testing.T) {
	t.Parallel()

	got, err := HashPassword(pw, &c.Secrets)
	if err != nil {
		t.Fatalf("bcrypt failed: %s", err)
	}

	err = VerifyPassword(got, pw)
	if err != nil {
		t.Errorf("%q should match %q", got, pw)
	}
}

func BenchmarkVerifyPassword(b *testing.B) {
	hash, err := HashPassword(pw, &c.Secrets)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		_ = VerifyPassword(hash, pw)
	}
}
