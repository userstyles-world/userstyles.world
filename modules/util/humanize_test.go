package util

import "testing"

var relNumberCases = []struct {
	name     string
	input    int64
	expected string
}{
	{"hundreds of millions", 777_777_777, "777.78M"},
	{"tens of millions", 77_777_777, "77.78M"},
	{"millions", 7_777_777, "7.78M"},
	{"hundreds of thousands", 777777, "777.78k"},
	{"tens of thousands", 42069, "42.07k"},
	{"thousands", 1337, "1.34k"},
	{"hundreds", 420, "420"},
	{"tens", 42, "42"},
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
