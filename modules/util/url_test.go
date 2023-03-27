package util

import "testing"

var slugCases = []struct {
	name, input, expected string
}{
	{"valid characters", "UserStyle Name", "userstyle-name"},
	{"invalid characters", "$(@#&($*#@%^#(@)))", "default-slug"},
	{"more characters", "What_Even-Is  This?!", "what-even-is-this"},
	{"extra characters", "(Dark) Theme [v1.2.3]", "dark-theme-v1-2-3"},
	{"uppercase characters", "MY USERSTYLE NAME", "my-userstyle-name"},
	{"many valid characters", "a b c d e f g h i", "a-b-c-d-e-f-g-h-i"},
	{"first typographic symbols", "暗い空 Dark Mode", "dark-mode"},
	{"first valid characters", "Dark Mode 暗い空", "dark-mode"},
	{"only typographic symbols", "暗い空", "default-slug"},
}

func TestSlug(t *testing.T) {
	t.Parallel()

	for _, c := range slugCases {
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

func BenchmarkSlug(b *testing.B) {
	for _, c := range slugCases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				Slug(c.input)
			}
		})
	}
}

var proxyCases = []struct {
	name     string
	link     string
	kind     string
	id       uint
	expected string
}{
	{
		"example text",
		"<h1>Hello, World!<h1>", "", 0,
		"<h1>Hello, World!<h1>",
	},
	{
		"example image",
		`<img src="https://example.com/foo.png">`, "style", 1,
		`<img src="/proxy?link=https://example.com/foo.png&type=style&id=1" loading="lazy">`,
	},
	{
		"oddly formatted example image",
		`<img  SRC  = " HTTP://EXAMPLE.COM/foo.png  " >`, "profile", 1,
		`<img  src="/proxy?link=HTTP://EXAMPLE.COM/foo.png&type=profile&id=1" loading="lazy" >`,
	},
	{
		"multiple example image",
		`<h1>hi</h1><img src="https://example.com/foo.png"><img src="https://example.com/bar.png">`, "style", 1,
		`<h1>hi</h1><img src="/proxy?link=https://example.com/foo.png&type=style&id=1" loading="lazy"><img src="/proxy?link=https://example.com/bar.png&type=style&id=1" loading="lazy">`,
	},
}

func TestProxyResources(t *testing.T) {
	t.Parallel()

	for _, c := range proxyCases {
		t.Run(c.name, func(t *testing.T) {
			got := ProxyResources(c.link, c.kind, c.id)
			if got != c.expected {
				t.Errorf("got: %s\n", got)
				t.Errorf("exp: %s\n", c.expected)
			}
		})
	}
}

func BenchmarkProxyResources(b *testing.B) {
	for _, c := range proxyCases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				ProxyResources(c.link, c.kind, c.id)
			}
		})
	}
}
