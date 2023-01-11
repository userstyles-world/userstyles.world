package util

import "testing"

func TestSlug(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name, input, expected string
	}{
		{"valid characters", "UserStyle Name", "userstyle-name"},
		{"invalid characters", "$(@#&($*#@%^#(@)))", "default-slug"},
		{"more characters", "What_Even-Is  This?!", "what-even-is-this"},
		{"extra characters", "(Dark) Theme [v1.2.3]", "dark-theme-v1-2-3"},
		{"many valid characters", "a b c d e f g h i", "a-b-c-d-e-f-g-h-i"},
		{"mixed typographic symbols", "暗い空 Dark Mode", "dark-mode"},
		{"only typographic symbols", "暗い空", "default-slug"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Slug(c.input)
			if got != c.expected {
				t.Errorf("got: %s\n", got)
				t.Errorf("exp: %s\n", c.expected)
			}
		})
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
