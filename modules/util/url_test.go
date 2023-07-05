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

var crawlerCases = []struct {
	name, input string
	expected    bool
}{
	{"random", "Some Random User Agent", false},
	{"firefox", "Mozilla/5.0 (X11; Linux x86_64; rv:109.0) Gecko/20100101 Firefox/112.0", false},
	{"chromium", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36", false},
	{"uppercase bot", "User Agent Of Some Random Bot V1.0", true},
	{"lowercase bot", "user agent of some random bot v1.0", true},
	{"lemmy", "Lemmy/0.17.2; +https://example.com", true},
	{"pict-rs", "pict-rs v0.3.0-main", true},
	{"calckey", "Calckey/14.0.0-dev42 (https://example.com)", true},
	{"misskey", "Misskey/13.12.2 (https://example.com", true},
	{"friendica", "Friendica 'Giant Rhubarb' 2023.05-1518; https://example.com", true},
	{"akkoma", "Akkoma 3.9.3-0-deadbeef; https://example.com <user@example.com>", true},
}

func TestIsCrawler(t *testing.T) {
	for _, c := range crawlerCases {
		t.Run(c.name, func(t *testing.T) {
			got := IsCrawler(c.input)
			if got != c.expected {
				t.Errorf("got: %t", got)
				t.Errorf("exp: %t", c.expected)
			}
		})
	}
}

func BenchmarkIsCrawler(b *testing.B) {
	for _, c := range crawlerCases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				IsCrawler(c.input)
			}
		})
	}
}
