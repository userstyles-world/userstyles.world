package strutils

import (
	"testing"
)

func TestSluggifyURLs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		desc     string
		a        string
		expected string
	}{
		{"TestProperName", "UserStyle Name", "userstyle-name"},
		{"TestMoreCharacters", "What_Even-Is  This?!", "what-even-is-this"},
		{"TestExtraCharacters", "(Dark) Something [v1.2.3]", "dark-something-v1-2-3"},
		{"TextExtraOfEverything", " Please---Get___Some   HELP!!! ", "please-get-some-help"},
		{"TestTypographicSymbols", "暗い空 Dark Mode", "dark-mode"},
		{"TestTypographicSymbolsOnly", "暗い空", "default-slug"},
	}

	for _, c := range cases {
		actual := SlugifyURL(c.a)
		if actual != c.expected {
			t.Fatalf("%s: expected: %s got: %s",
				c.desc, c.expected, actual)
		}
	}
}

func TestHumanizeNumber(t *testing.T) {
	t.Parallel()

	cases := []struct {
		desc     string
		arg      int
		expected string
	}{
		{"Hundreds of thousands", 777777, "777.7k"},
		{"Tens of thousands", 42069, "42k"},
		{"Thousands", 1337, "1.3k"},
		{"Hundreds", 420, "420"},
		{"Tens", 42, "42"},
	}

	for _, c := range cases {
		got := HumanizeNumber(c.arg)
		if got != c.expected {
			t.Fatalf("%s: expected: %s got: %s", c.desc, c.expected, got)
		}
	}
}
