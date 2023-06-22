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
	{"100T", 100e12, "100T"},
	{"10T", 10e12, "10T"},
	{"1T", 1e12 - 49999, "1T"},
	{"100B", 100e9, "100B"},
	{"10B", 10e9, "10B"},
	{"1B", 1e9 - 4999, "1B"},
	{"100M", 100e6, "100M"},
	{"10M", 10e6, "10M"},
	{"1M", 999501, "1M"},
	{"999.5k", 999500, "999.5k"},
	{"420.4k", 420420, "420.4k"},
	{"111k", 110951, "111k"},
	{"100k", 100000, "100k"},
	{"99.9k", 99950, "99.9k"},
	{"42.1k", 42069, "42.1k"},
	{"10.9k", 10950, "10.9k"},
	{"10k", 10000, "10k"},
	{"1k", 1337, "1337"},
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
	{"1s", time.Now().Add(-1 * second), "1 second ago"},
	{"2s", time.Now().Add(-2 * second), "2 seconds ago"},
	{"1m", time.Now().Add(-1 * minute), "1 minute ago"},
	{"2m", time.Now().Add(-2 * minute), "2 minutes ago"},
	{"1h", time.Now().Add(-1 * hour), "1 hour ago"},
	{"2h", time.Now().Add(-2 * hour), "2 hours ago"},
	{"1d", time.Now().Add(-1 * day), "1 day ago"},
	{"2d", time.Now().Add(-2 * day), "2 days ago"},
	{"1w", time.Now().Add(-1 * week), "1 week ago"},
	{"2w", time.Now().Add(-2 * week), "2 weeks ago"},
	{"1mo", time.Now().Add(-month), "1 month ago"},
	{"2mo", time.Now().Add(-2 * month), "2 months ago"},
	{"1y", time.Now().Add(-1 * year), "1 year ago"},
	{"2y", time.Now().Add(-2 * year), "2 years ago"},
	{"1m9s", time.Now().Add(-69 * second), "1 minute, 9 seconds ago"},
	{"2h1m", time.Now().Add(-121 * minute), "2 hours, 1 minute ago"},
	{"2h46m", time.Now().Add(-9999 * second), "2 hours, 46 minutes ago"},
	{"6d22h", time.Now().Add(-9999 * minute), "6 days, 22 hours ago"},
	{"1w3d", time.Now().Add(-10 * day), "1 week, 3 days ago"},
	{"4w2d", time.Now().Add(-30 * day), "4 weeks, 2 days ago"},
	{"1m1d", time.Now().Add(-32 * day), "1 month, 1 day ago"},
	{"3m1w", time.Now().Add(-99 * day), "3 months, 1 week ago"},
	{"1y18h", time.Now().Add(-366 * day), "1 year, 18 hours ago"},
	{"3y5w", time.Now().Add(-42 * month), "3 years, 5 months ago"},
	{"34y11mo", time.Now().Add(-420 * month), "34 years, 11 months ago"},
	{"11mo4w", time.Now().Add(-52 * week), "11 months, 4 weeks ago"},
	{"future", time.Now().Add(420 * time.Hour), "2 weeks, 3 days"},
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
