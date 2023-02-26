package util

import (
	"fmt"
	"testing"
)

const ucTpl = `/* ==UserStyle==
@name         test
@namespace    test
@version      1.0.0
%s
==UserStyle== */`

const ucExp = `/* ==UserStyle==
@name         test
@namespace    test
@version      1.0.0
==UserStyle== */`

func TestRemoveUpdateURL(t *testing.T) {
	t.Parallel()

	gen := func(s string) string {
		return fmt.Sprintf(ucTpl, s)
	}

	cases := []struct {
		name, input, expected string
	}{
		{"normal", gen("@updateURL https://example.org/1.user.css"), ucExp},
		{"spaces", gen(" @updateURL https://example.org/1.user.css "), ucExp},
		{"tabs", gen("	@updateURL https://example.org/1.user.css	"), ucExp},
		{"both", gen(" 	@updateURL https://example.org/1.user.css 	"), ucExp},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := RemoveUpdateURL(c.input)
			if got != c.expected {
				t.Errorf("got: %v\n", got)
				t.Errorf("exp: %v\n", c.expected)
			}
		})
	}
}
