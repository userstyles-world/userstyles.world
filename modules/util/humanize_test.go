package util

import "testing"

func TestRelNumber(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    int
		expected string
	}{
		{"Hundreds of thousands", 777777, "777.7k"},
		{"Tens of thousands", 42069, "42k"},
		{"Thousands", 1337, "1.3k"},
		{"Hundreds", 420, "420"},
		{"Tens", 42, "42"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := RelNumber(c.input)
			if got != c.expected {
				t.Errorf("got: %s\n", got)
				t.Errorf("exp: %s\n", c.expected)
			}
		})
	}
}
