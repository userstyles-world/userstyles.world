package images

import (
	"testing"
)

func TestFixRawURL(t *testing.T) {
	t.Parallel()

	cases := []struct {
		desc     string
		a        string
		expected string
	}{
		{"TestGitHubBlobURL", "https://github.com/user/style/blob/main/img.png", "https://github.com/user/style/raw/main/img.png"},
		{"TestGitHubRawURL", "https://github.com/user/style/raw/main/img.png", "https://github.com/user/style/raw/main/img.png"},
		{"TestGitLabBlobURL", "https://gitlab.com/user/style/-/blob/main/img.png", "https://gitlab.com/user/style/-/raw/main/img.png"},
		{"TestGitLabRawURL", "https://gitlab.com/user/style/-/raw/main/img.png?inline=false", "https://gitlab.com/user/style/-/raw/main/img.png"},
		{"TestCodebergBlobURL", "https://codeberg.org/user/style/src/branch/main/img.png", "https://codeberg.org/user/style/raw/branch/main/img.png"},
		{"TestCodebergRawURL", "https://codeberg.org/user/style/raw/branch/main/img.png", "https://codeberg.org/user/style/raw/branch/main/img.png"},
	}

	for _, c := range cases {
		actual := fixRawURL(c.a)
		if actual != c.expected {
			t.Fatalf("%s: expected: %s got: %s", c.desc, c.expected, actual)
		}
	}
}
