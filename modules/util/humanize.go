package util

import (
	"strconv"
	"sync"
	"unsafe"
)

const nSmalls = 100 // taken from strconv/itoa.go

var relNumberPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 0, 10)
		return &buf
	},
}

// RelNumber returns a relative representation of a number i.
func RelNumber(i int64) string {
	if i < nSmalls {
		return strconv.FormatInt(i, 10)
	}

	bp := relNumberPool.Get().(*[]byte)
	defer relNumberPool.Put(bp)
	b := (*bp)[:0]

	switch {
	case i > 1_000_000:
		b = strconv.AppendFloat(b, float64(i)/1_000_000, 'f', 1, 64)
		b = append(b, 'M')
	case i >= 1_000:
		b = strconv.AppendFloat(b, float64(i)/1_000, 'f', 1, 64)
		b = append(b, 'k')
	default:
		b = strconv.AppendInt(b, i, 10)
	}

	return *(*string)(unsafe.Pointer(&b))
}
