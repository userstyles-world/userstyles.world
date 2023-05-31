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
			name: "mirrored import public",
			input: APIStyle{
				Original:   "x",
				MirrorCode: true,
			},
			branch:   true,
			combined: "Imported and mirrored from <code>x</code>",
		},
		{
			name: "mirrored import private due to import",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: true,
				MirrorCode:    true,
			},
			branch:   true,
			combined: "Imported and mirrored from a private source",
		},
		{
			name: "mirrored import private due to mirror",
			input: APIStyle{
				Original:      "x",
				MirrorPrivate: true,
				MirrorCode:    true,
			},
			branch:   true,
			combined: "Imported and mirrored from a private source",
		},
		{
			name: "both public",
			input: APIStyle{
				Original:   "x",
				MirrorURL:  "x",
				MirrorCode: true,
			},
			branch:   true,
			combined: "Imported and mirrored from <code>x</code>",
		},
		{
			name: "both private due to import",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: true,
				MirrorURL:     "x",
				MirrorCode:    true,
			},
			branch:   true,
			combined: "Imported and mirrored from a private source",
		},
		{
			name: "both private due to mirror",
			input: APIStyle{
				Original:      "x",
				MirrorURL:     "x",
				MirrorPrivate: true,
				MirrorCode:    true,
			},
			branch:   true,
			combined: "Imported and mirrored from a private source",
		},
		{
			name: "different URLs public",
			input: APIStyle{
				Original:   "x",
				MirrorURL:  "y",
				MirrorCode: true,
			},
			branch:   false,
			imported: "Imported from <code>x</code>",
			mirrored: "Mirrored from <code>y</code>",
		},
		{
			name: "different URLs private",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: true,
				MirrorURL:     "y",
				MirrorPrivate: true,
				MirrorCode:    true,
			},
			branch:   false,
			imported: "Imported from a private source",
			mirrored: "Mirrored from a private source",
		},
		{
			name: "different URLs private import",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: true,
				MirrorURL:     "y",
				MirrorCode:    true,
			},
			branch:   false,
			imported: "Imported from a private source",
			mirrored: "Mirrored from <code>y</code>",
		},
		{
			name: "different URLs private mirror",
			input: APIStyle{
				Original:      "x",
				MirrorURL:     "y",
				MirrorPrivate: true,
				MirrorCode:    true,
			},
			branch:   false,
			imported: "Imported from <code>x</code>",
			mirrored: "Mirrored from a private source",
		},
		{
			name: "has import URL but not mirrored",
			input: APIStyle{
				Original: "x",
			},
			branch:   false,
			imported: "Imported from <code>x</code>",
		},
		{
			name: "has mirror URL but not mirrored",
			input: APIStyle{
				MirrorURL: "y",
			},
			branch: false,
		},
		{
			name:  "empty",
			input: APIStyle{},
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
				if c.input.Imported() {
					got := c.input.ImportedText()
					if got != c.imported {
						t.Errorf("got: %v\n", got)
						t.Errorf("exp: %v\n", c.imported)
					}
				}
				if c.input.Mirrored() {
					got := c.input.MirroredText()
					if got != c.mirrored {
						t.Errorf("got: %v\n", got)
						t.Errorf("exp: %v\n", c.mirrored)
					}
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
			name: "public",
			input: APIStyle{
				Original: "x",
			},
			expected: "Imported from <code>x</code>",
		},
		{
			name: "private",
			input: APIStyle{
				Original:      "x",
				ImportPrivate: true,
			},
			expected: "Imported from a private source",
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
			name: "public",
			input: APIStyle{
				MirrorURL:  "x",
				MirrorCode: true,
			},
			expected: "Mirrored from <code>x</code>",
		},
		{
			name: "private",
			input: APIStyle{
				MirrorURL:     "x",
				MirrorPrivate: true,
				MirrorCode:    true,
			},
			expected: "Mirrored from a private source",
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
