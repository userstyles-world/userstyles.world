package util

import (
	"testing"
)

var hashCases = []struct {
	name, id, ip, expected string
}{
	{"foobar", "foo", "bar", "4c7512b2468ccab5d8e3cfc3ee7ec0c76f154c1d145c23f716748fd3d6789e0b2283b46fcbcc660134a54a8e71d5171fc16df69b6c64896639eb4c9f49cf035c"},
	{"real IP", "1", "192.168.1.1", "befd488b1e46e011eeeb0078cb7f6f10792068e5105bcf98ab0e2a7ed31c3863695b84ed23dc672c10dad4c815b37d418bfc3640a1a47eddcdde0c352f472f00"},
	{"no IP", "1", "", "b7d09939dd47d20d5460ca664cc08919f8c89bfca3233ea2bdfadfa32efde6bad2e4261f6f0542171bc4e958add12d29605f2d97a901006592a34ce8a5a0cd42"},
	{"no ID", "", "192.168.1.1", "1778721963e12865be624b66e70103c14cbb1c2ccc494a990927edb89f67eb58232ac8bad231c75192643808b75a6d86228b5b03da2743af748c0cc1be68aa9f"},
}

func TestHashIP(t *testing.T) {
	t.Parallel()
	InitPool(&c.Secrets)

	for _, c := range hashCases {
		actual, err := HashIP(c.ip + " " + c.id)
		if err != nil {
			t.Error(err)
		}
		if actual != c.expected {
			t.Fatalf("%s: expected: %s got: %s",
				c.name, c.expected, actual)
		}
	}
}

func BenchmarkHashIP(b *testing.B) {
	for _, c := range hashCases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = HashIP(c.ip + " " + c.id)
			}
		})
	}
}
