package util

import (
	"fmt"
	"regexp"
	"strconv"
	"sync"
	"unicode"
	"unsafe"
)

var (
	linkRe   = regexp.MustCompile(`(?mU)src="(http.*)"`)
	slugPool = sync.Pool{
		New: func() any {
			buf := make([]byte, 0, 256)
			return &buf
		},
	}
)

// Slug takes in a string s and returns a user- and SEO-friendly URL.
func Slug(s string) string {
	bp := slugPool.Get().(*[]byte)
	defer slugPool.Put(bp)
	b := (*bp)[:0]

	var sep bool
	for _, c := range s {
		switch {
		case (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'):
			b = append(b, byte(unicode.ToLower(c)))
			sep = true
		case c == ' ' || c == '-' || c == '_' || c == '.':
			if sep {
				b = append(b, '-')
				sep = false
			}
		}
	}

	if len(b) == 0 {
		return "default-slug"
	}

	return *(*string)(unsafe.Pointer(&b))
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
