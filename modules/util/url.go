package util

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
	"unsafe"
)

var (
	linkRe   = regexp.MustCompile(`(?imU)src\s*=\s*['"]\s*(http.*)\s*['"]`)
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
		case c >= 'A' && c <= 'Z':
			b = append(b, byte(c|32)) // [:lower:] = [:upper:] | 32
			sep = true
		case (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9'):
			b = append(b, byte(c))
			sep = true
		case (c == ' ' || c == '-' || c == '_' || c == '.') && sep:
			b = append(b, '-')
			sep = false
		}
	}

	if len(b) == 0 {
		return "default-slug"
	}

	// Trim trailing dash.
	if b[len(b)-1] == '-' {
		b = b[:len(b)-1]
	}

	return *(*string)(unsafe.Pointer(&b))
}

// ProxyResources takes in external images and stores them locally.
func ProxyResources(s, t string, id uint) string {
	sub := fmt.Sprintf(`src="/proxy?link=$1&type=%s&id=%d" loading="lazy"`, t, id)
	return linkRe.ReplaceAllString(s, sub)
}

// IsCrawler ignores crawlers in places where we collect statistics.
func IsCrawler(ua string) bool {
	ua = strings.ToLower(ua)
	return strings.Contains(ua, "bot")
}
