package images

import (
	"regexp"
)

func fixRawURL(url string) string {
	re := regexp.MustCompile(`(?mi)^(http.*)/(raw|src|blob)/(.*.(png|jpe?g|avif|webp))(\?.*)*$`)
	return re.ReplaceAllString(url, "${1}/raw/${3}")
}
