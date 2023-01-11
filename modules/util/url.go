package util

import (
	"fmt"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var (
	slugRe = regexp.MustCompile(`[a-zA-Z0-9]+`)
	linkRe = regexp.MustCompile(`(?mU)src="(http.*)"`)
)

// Slug takes in a string s and returns a user- and SEO-friendly URL.
func Slug(s string) string {
	// Extract valid characters.
	parts := slugRe.FindAllString(s, -1)

	// Add default slug for unsupported characters.
	if len(parts) == 0 {
		return "default-slug"
	}

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

func HumanizeNumber(i int) string {
	switch {
	case i >= 100_000:
		return format(i, 3)

	case i >= 10_000:
		return format(i, 2)

	case i >= 1_000:
		return format(i, 1)

	default:
		return strconv.Itoa(i)
	}
}

func format(i, pos int) string {
	s := fmt.Sprintf("%d", i)

	if s[pos] == '0' {
		return fmt.Sprintf("%sk", s[:pos])
	}

	return fmt.Sprintf("%s.%ck", s[:pos], s[pos])
}
