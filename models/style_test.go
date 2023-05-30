package models

import "testing"

func TestAPIStyle_ImportedAndMirrored(t *testing.T) {
	cases := []struct {
		name     string
		input    APIStyle
		expected string
	}{
		{
			name: "private import",
			input: APIStyle{
				Original:      "x",
				MirrorCode:    true,
				ImportPrivate: true,
			},
			expected: "Imported and mirrored from a private source",
		},
		{
			name: "public origin",
			input: APIStyle{
				Original:   "x",
				MirrorCode: true,
			},
			expected: "Imported and mirrored from <code>x</code>",
		},
		{
			name: "both private",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: true,
				MirrorURL:     "x",
				MirrorCode:    true,
				MirrorPrivate: false,
			},
			expected: "Imported and mirrored from a private source",
		},
		{
			name: "both public",
			input: APIStyle{
				Original:   "x",
				MirrorCode: true,
				MirrorURL:  "x",
			},
			expected: "Imported and mirrored from <code>x</code>",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if !c.input.ImportedAndMirrored() {
				t.Fatal("import and mirror URL don't match")
			}
			if got := c.input.ImportedAndMirroredText(); got != c.expected {
				t.Errorf("got: %v\n", got)
				t.Errorf("exp: %v\n", c.expected)
			}
		})
	}
}

func TestAPIStyle_Imported(t *testing.T) {
	cases := []struct {
		name     string
		input    APIStyle
		expected string
	}{
		{
			name: "private",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: true,
			},
			expected: "Imported from a private source",
		},
		{
			name: "public",
			input: APIStyle{
				Original: "x",
			},
			expected: "Imported from <code>x</code>",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if !c.input.Imported() {
				t.Errorf("style isn't imported")
			}
			if got := c.input.ImportedText(); got != c.expected {
				t.Errorf("got: %v\n", got)
				t.Errorf("exp: %v\n", c.expected)
			}
		})
	}
}

func TestAPIStyle_Mirrored(t *testing.T) {
	cases := []struct {
		name     string
		input    APIStyle
		expected string
	}{
		{
			name: "private",
			input: APIStyle{
				MirrorURL:     "x",
				MirrorCode:    true,
				MirrorPrivate: true,
			},
			expected: "Mirrored from a private source",
		},
		{
			name: "public",
			input: APIStyle{
				MirrorURL:  "x",
				MirrorCode: true,
			},
			expected: "Mirrored from <code>x</code>",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if !c.input.Mirrored() {
				t.Errorf("style isn't mirrored")
			}
			if got := c.input.MirroredText(); got != c.expected {
				t.Errorf("got: %v\n", got)
				t.Errorf("exp: %v\n", c.expected)
			}
		})
	}
}
