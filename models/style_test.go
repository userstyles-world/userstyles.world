package models

import "testing"

var importAndMirrorCases = []struct {
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

func TestAPIStyle_ImportedAndMirrored(t *testing.T) {
	for _, c := range importAndMirrorCases {
		t.Run(c.name, func(t *testing.T) {
			if c.branch {
				combined := c.input.ImportedAndMirrored()
				if combined != c.combined {
					t.Errorf("combined: %v\n", combined)
					t.Errorf("expected: %v\n", c.combined)
				}
			} else {
				if c.input.isImportedAndMirrored() {
					t.Fatal("import and mirror URL should match")
				}
				imported := c.input.Imported()
				if imported != c.imported {
					t.Errorf("imported: %v\n", imported)
					t.Errorf("expected: %v\n", c.imported)
				}
				mirrored := c.input.Mirrored()
				if mirrored != c.mirrored {
					t.Errorf("mirrored: %v\n", mirrored)
					t.Errorf("expected: %v\n", c.mirrored)
				}
			}
		})
	}
}

func BenchmarkAPIStyle_ImportedAndMirrored(b *testing.B) {
	for _, c := range importAndMirrorCases {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if c.branch {
					c.input.ImportedAndMirrored()
				} else {
					c.input.Imported()
					c.input.Mirrored()
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
			if !c.input.isImported() {
				t.Fatal("style isn't imported")
			}
			if got := c.input.Imported(); got != c.expected {
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
			if !c.input.isMirrored() {
				t.Fatal("style isn't mirrored")
			}
			if got := c.input.Mirrored(); got != c.expected {
				t.Errorf("got: %v\n", got)
				t.Errorf("exp: %v\n", c.expected)
			}
		})
	}
}
