package util

import "regexp"

var urlRe = regexp.MustCompile(`(?m)^\s*@updateURL.*\n`)

// RemoveUpdateURL strips away `@updateURL` field from userstyles.
func RemoveUpdateURL(s string) string {
	return urlRe.ReplaceAllString(s, "")
}
