package utils

import (
	"testing"
)

func BenchmarkRandomNonce(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandomBytes(24)
	}
}

func BenchmarkRandomBoundary(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandomString(30)
	}
}

func BenchmarkRandomOauth(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandomString(128)
	}
}

func TestRandomString(t *testing.T) {
	t.Parallel()
	for i := 0; i < 100; i++ {
		s := RandomString(20)
		if len(s) != 40 {
			t.Error("RandomString returned wrong length")
		}
	}
}
