package utils

import (
	"testing"

	"userstyles.world/utils"
)

func BenchmarkRandomNonce(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.RandStringBytesMaskImprSrcUnsafe(24)
	}
}

func BenchmarkRandomBoundary(b *testing.B) {
	for n := 0; n < b.N; n++ {
		utils.RandStringBytesMaskImprSrcUnsafe(30)
	}
}
