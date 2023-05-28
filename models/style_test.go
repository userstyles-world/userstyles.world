package models

import "testing"

func TestStyle(t *testing.T) {
	cases := []struct {
		name     string
		input    APIStyle
		expected string
	}{
		{
			name: "imported",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: true,
			},
			expected: "<p class=\"mb:m md\">Imported</p>",
		},
		{
			name: "imported from x",
			input: APIStyle{
				Original: "x",
			},
			expected: "<p class=\"mb:m md\">Imported from <code>x</code></p>",
		},
		{
			name: "mirrored",
			input: APIStyle{
				MirrorURL:     "x",
				MirrorCode:    true,
				MirrorPrivate: true,
			},
			expected: "<p class=\"mb:m md\">Mirrored</p>",
		},
		{
			name: "mirrored from x",
			input: APIStyle{
				MirrorURL:  "x",
				MirrorCode: true,
			},
			expected: "<p class=\"mb:m md\">Mirrored from <code>x</code></p>",
		},
		{
			name: "imported and mirrored",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: true,
				MirrorURL:     "y",
				MirrorCode:    true,
				MirrorPrivate: true,
			},
			expected: "<p class=\"mb:m md\">Imported</p><p class=\"mb:m md\">Mirrored</p>",
		},
		{
			name: "imported from x and mirrored",
			input: APIStyle{
				Original:      "x",
				MirrorURL:     "y",
				MirrorCode:    true,
				MirrorPrivate: true,
			},
			expected: "<p class=\"mb:m md\">Imported from <code>x</code></p><p class=\"mb:m md\">Mirrored</p>",
		},
		{
			name: "imported and mirrored from x",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: true,
				MirrorURL:     "y",
				MirrorCode:    true,
			},
			expected: "<p class=\"mb:m md\">Imported</p><p class=\"mb:m md\">Mirrored from <code>y</code></p>",
		},
		{
			name: "imported from x and mirrored from x",
			input: APIStyle{
				Original:   "x",
				MirrorURL:  "y",
				MirrorCode: true,
			},
			expected: "<p class=\"mb:m md\">Imported from <code>x</code></p><p class=\"mb:m md\">Mirrored from <code>y</code></p>",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.input.HeadlineText(); got != c.expected {
				t.Errorf("got: %s\n", got)
				t.Errorf("exp: %s\n", c.expected)
			}
		})
	}
}
