package util

import "testing"

var relNumberCases = []struct {
	name     string
	input    int64
	expected string
}{
	{"100M", 100e6, "100.00M"},
	{"10M", 10e6, "10.00M"},
	{"1M", 1e6, "1.00M"},
	{"100k", 100e3, "100.00k"},
	{"10k", 42069, "42.07k"},
	{"1k", 1337, "1.34k"},
	{"100", 420, "420"},
	{"10", 42, "42"},
	{"0", 0, "0"},
}

func TestRelNumber(t *testing.T) {
	t.Parallel()

	for _, c := range relNumberCases {
		t.Run(c.name, func(t *testing.T) {
			got := RelNumber(c.input)
			if got != c.expected {
				t.Errorf("got: %s\n", got)
				t.Errorf("exp: %s\n", c.expected)
			}
		})
	}
}

func BenchmarkRelNumber(b *testing.B) {
	for _, c := range relNumberCases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				RelNumber(c.input)
			}
		})
	}
}
