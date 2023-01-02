package util

import (
	"fmt"
	"strconv"
)

// RelNumber returns a relative representation of a number i.
func RelNumber(i int) string {
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
