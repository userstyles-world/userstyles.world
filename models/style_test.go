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
			expected: "Imported.",
		},
		{
			name: "imported from x",
			input: APIStyle{
				Original: "x",
			},
			expected: "Imported from <code>x</code>.",
		},
		{
			name: "mirrored",
			input: APIStyle{
				MirrorURL:     "x",
				MirrorCode:    true,
				MirrorPrivate: true,
			},
			expected: "Mirrored.",
		},
		{
			name: "mirrored from x",
			input: APIStyle{
				MirrorURL:  "x",
				MirrorCode: true,
			},
			expected: "<nobr>Mirrored from <code>x</code>.</nobr>",
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
			expected: "Imported, Mirrored.",
		},
		{
			name: "imported from x and mirrored",
			input: APIStyle{
				Original:      "x",
				MirrorURL:     "y",
				MirrorCode:    true,
				MirrorPrivate: true,
			},
			expected: "Imported from <code>x</code>, Mirrored.",
		},
		{
			name: "imported and mirrored from x",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: true,
				MirrorURL:     "y",
				MirrorCode:    true,
			},
			expected: "Imported, Mirrored from <code>y</code>.",
		},
		{
			name: "imported from x and mirrored from x",
			input: APIStyle{
				Original:   "x",
				MirrorURL:  "y",
				MirrorCode: true,
			},
			expected: "Imported from <code>x</code>, <nobr>Mirrored from <code>y</code>.</nobr>",
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
