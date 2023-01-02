package util

import (
	"strconv"
)

// RelNumber returns a relative representation of a number i.
func RelNumber(i int) string {
	switch {
	case i >= 1_000:
		return strconv.FormatFloat(float64(i)/1_000, 'f', 1, 64) + "k"
	default:
		return strconv.Itoa(i)
	}
}
