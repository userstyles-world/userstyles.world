package util

import (
	"strconv"
	"sync"
	"time"
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
	case i >= 1e12-5e6:
		b = strconv.AppendFloat(b, float64(i)/1e12, 'f', 2, 32)
		b = append(b, 'T')
	case i >= 1e9-5e3:
		b = strconv.AppendFloat(b, float64(i)/1e9, 'f', 2, 32)
		b = append(b, 'B')
	case i >= 1e6-5:
		b = strconv.AppendFloat(b, float64(i)/1e6, 'f', 2, 32)
		b = append(b, 'M')
	case i >= 1e3:
		b = strconv.AppendFloat(b, float64(i)/1e3, 'f', 2, 32)
		b = append(b, 'K')
	default:
		b = strconv.AppendInt(b, i, 10)
	}

	return *(*string)(unsafe.Pointer(&b))
}

func buildTime(b []byte, parts *int, count int, dur string) []byte {
	*parts++
	b = strconv.AppendInt(b, int64(count), 10)
	b = append(b, dur...)
	if count > 1 {
		b = append(b, "s, "...)
	} else {
		b = append(b, ", "...)
	}

	return b
}

// RelTime returns a relative representation of a time t.
func RelTime(t time.Time) string {
	var (
		parts   = 0
		b       = []byte{}
		now     = time.Since(t)
		hours   = int(now / time.Hour % 60)
		minutes = int(now / time.Minute % 60)
		seconds = int(now / time.Second % 60)
	)

	if hours > 0 && parts != 2 {
		b = buildTime(b, &parts, hours, " hour")
	}
	if minutes > 0 && parts != 2 {
		b = buildTime(b, &parts, minutes, " minute")
	}
	if seconds > 0 && parts != 2 {
		b = buildTime(b, &parts, seconds, " second")
	}

	if len(b) == 0 {
		return "just now"
	}

	b = append(b[:len(b)-2], " ago"...)

	return string(b)
}
