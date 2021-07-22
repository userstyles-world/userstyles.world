package utils

import (
	"strings"
	"testing"
)

func BenchmarkEqual(b *testing.B) {
	b.StopTimer()
	passwd := "somepasswordyoulike"
	hash := GenerateHashedPassword(passwd)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = CompareHashedPassword(hash, passwd)
	}
}

func BenchmarkDefaultCost(b *testing.B) {
	b.StopTimer()
	passwd := "mylongpassword1234"
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		GenerateHashedPassword(passwd)
	}
}

func TestBcryptingIsCorrect(t *testing.T) {
	pass := "allmine"
	expectedHash := "$2a$10"

	hash := GenerateHashedPassword(pass)
	if !strings.HasPrefix(hash, expectedHash) {
		t.Errorf("%v should be the suffix of %v", hash, expectedHash)
	}
}
