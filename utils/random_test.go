package utils

import "testing"

func BenchmarkRandomNonce(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		RandStringBytesMaskImprSrcUnsafe(12)
	}
}

func BenchmarkRandomBoundary(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		RandStringBytesMaskImprSrcUnsafe(30)
	}
}
