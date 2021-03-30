package utils

import "testing"

func BenchmarkRandomNonce(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandStringBytesMaskImprSrcUnsafe(12)
	}
}

func BenchmarkRandomBoundary(b *testing.B) {
	for n := 0; n < b.N; n++ {
		RandStringBytesMaskImprSrcUnsafe(30)
	}
}
