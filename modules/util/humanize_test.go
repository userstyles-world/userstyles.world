package util

import (
	"testing"
	"time"

	"github.com/dustin/go-humanize"
)

var relNumberCases = []struct {
	name     string
	input    int64
	expected string
}{
	{"100T", 100e12, "100.00T"},
	{"10T", 10e12, "10.00T"},
	{"1T", 1e12 - 5e6, "1.00T"},
	{"100B", 100e9, "100.00B"},
	{"10B", 10e9, "10.00B"},
	{"1B", 1e9 - 5e3, "1.00B"},
	{"100M", 100e6, "100.00M"},
	{"10M", 10e6, "10.00M"},
	{"1M", 1e6 - 5, "1.00M"},
	{"100K", 100e3, "100.00K"},
	{"10K", 42069, "42.07K"},
	{"1K", 1337, "1.34K"},
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

var relTimeCases = []struct {
	name     string
	input    time.Time
	expected string
}{
	{"now", time.Now(), "just now"},
	{"1s", time.Now().Add(-1 * time.Second), "1 second ago"},
	{"2s", time.Now().Add(-2 * time.Second), "2 seconds ago"},
	{"1m", time.Now().Add(-1 * time.Minute), "1 minute ago"},
	{"2m", time.Now().Add(-2 * time.Minute), "2 minutes ago"},
	{"1h", time.Now().Add(-1 * time.Hour), "1 hour ago"},
	{"2h", time.Now().Add(-2 * time.Hour), "2 hours ago"},
	{"1m9s", time.Now().Add(-69 * time.Second), "1 minute, 9 seconds ago"},
	{"2h1m", time.Now().Add(-121 * time.Minute), "2 hours, 1 minute ago"},
	{"2h1m9s", time.Now().Add(-2 * time.Hour).Add(-time.Minute).Add(-9 * time.Second), "2 hours, 1 minute, 9 seconds ago"},
}

func TestRelTime(t *testing.T) {
	t.Parallel()

	for _, c := range relTimeCases {
		t.Run(c.name, func(t *testing.T) {
			got := RelTime(c.input)
			if got != c.expected {
				t.Errorf("got: %s\n", got)
				t.Errorf("exp: %s\n", c.expected)
			}
		})
	}
}

func BenchmarkRelTime(b *testing.B) {
	for _, c := range relTimeCases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				RelTime(c.input)
			}
		})
	}
}

func BenchmarkHumanizeTime(b *testing.B) {
	for _, c := range relTimeCases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				humanize.Time(c.input)
			}
		})
	}
}
