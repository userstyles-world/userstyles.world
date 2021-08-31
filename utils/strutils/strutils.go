package strutils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	slugRe = regexp.MustCompile(`[a-zA-Z0-9]+`)
	linkRe = regexp.MustCompile(`(?mU)src="(http.*)"`)
)

func SlugifyURL(s string) string {
	// Extract valid characters.
	parts := slugRe.FindAllString(s, -1)

	// Join parts and make them lowercase.
	s = strings.Join(parts, "-")
	s = strings.ToLower(s)

	return s
}

func QueryUnescape(s string) string {
	s, err := url.QueryUnescape(s)
	if err != nil {
		s = err.Error()
	}

	return s
}
func ProxyResources(s, t string, id uint) string {
	sub := fmt.Sprintf(`src="/proxy?link=$1&type=%s&id=%d" loading="lazy"`, t, id)
	return linkRe.ReplaceAllString(s, sub)
}
