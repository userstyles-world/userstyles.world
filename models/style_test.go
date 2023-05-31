package models

import "testing"

func TestAPIStyle_ImportedAndMirrored(t *testing.T) {
	cases := []struct {
		name     string
		input    APIStyle
		branch   bool
		combined string
		imported string
		mirrored string
	}{
		{
			name: "private import",
			input: APIStyle{
				Original:      "x",
				MirrorCode:    true,
				ImportPrivate: true,
			},
			branch:   true,
			combined: "Imported and mirrored from a private source",
		},
		{
			name: "public origin",
			input: APIStyle{
				Original:   "x",
				MirrorCode: true,
			},
			branch:   true,
			combined: "Imported and mirrored from <code>x</code>",
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
			branch:   true,
			combined: "Imported and mirrored from a private source",
		},
		{
			name: "both public",
			input: APIStyle{
				Original:   "x",
				MirrorCode: true,
				MirrorURL:  "x",
			},
			branch:   true,
			combined: "Imported and mirrored from <code>x</code>",
		},
		{
			name: "different URLs private",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: true,
				MirrorCode:    true,
				MirrorURL:     "y",
				MirrorPrivate: true,
			},
			branch:   false,
			imported: "Imported from a private source",
			mirrored: "Mirrored from a private source",
		},
		{
			name: "different URLs public",
			input: APIStyle{
				Original:   "x",
				MirrorCode: true,
				MirrorURL:  "y",
			},
			branch:   false,
			imported: "Imported from <code>x</code>",
			mirrored: "Mirrored from <code>y</code>",
		},
		{
			name: "different URLs private import",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: true,
				MirrorCode:    true,
				MirrorURL:     "y",
				MirrorPrivate: false,
			},
			branch:   false,
			imported: "Imported from a private source",
			mirrored: "Mirrored from <code>y</code>",
		},
		{
			name: "different URLs private mirror",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: false,
				MirrorCode:    true,
				MirrorURL:     "y",
				MirrorPrivate: true,
			},
			branch:   false,
			imported: "Imported from <code>x</code>",
			mirrored: "Mirrored from a private source",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if c.input.ImportedAndMirrored() != c.branch {
				t.Fatal("import and mirror URL should match")
			}

			if c.branch {
				got := c.input.ImportedAndMirroredText()
				if got != c.combined {
					t.Errorf("got: %v\n", got)
					t.Errorf("exp: %v\n", c.combined)
				}
			} else {
				got := c.input.ImportedText()
				if got != c.imported {
					t.Errorf("got: %v\n", got)
					t.Errorf("exp: %v\n", c.imported)
				}
				got = c.input.MirroredText()
				if got != c.mirrored {
					t.Errorf("got: %v\n", got)
					t.Errorf("exp: %v\n", c.mirrored)
				}
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
				t.Fatal("style isn't imported")
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
				t.Fatal("style isn't mirrored")
			}
			if got := c.input.MirroredText(); got != c.expected {
				t.Errorf("got: %v\n", got)
				t.Errorf("exp: %v\n", c.expected)
			}
		})
	}
}
