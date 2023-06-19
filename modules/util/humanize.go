package util

import (
	"bytes"
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
	case i >= 1e12-5e6+1:
		b = strconv.AppendFloat(b, float64(i)/1e12, 'f', 2, 32)
		b = bytes.TrimSuffix(b, []byte(".00"))
		b = append(b, 'T')
	case i >= 1e9-5e3+1:
		b = strconv.AppendFloat(b, float64(i)/1e9, 'f', 2, 32)
		b = bytes.TrimSuffix(b, []byte(".00"))
		b = append(b, 'B')
	case i >= 1e6-5e2+1:
		b = strconv.AppendFloat(b, float64(i)/1e6, 'f', 2, 32)
		b = bytes.TrimSuffix(b, []byte(".00"))
		b = append(b, 'M')
	case i >= 1e4:
		b = strconv.AppendFloat(b, float64(i)/1e3, 'f', 1, 32)
		b = bytes.TrimSuffix(b, []byte(".0"))
		b = append(b, 'K')
	default:
		b = strconv.AppendInt(b, i, 10)
	}

	return *(*string)(unsafe.Pointer(&b))
}

const (
	maxParts = 2
	second   = 1e9
	minute   = 60 * second
	hour     = 60 * minute
	day      = 24 * hour
	week     = 7 * day
	month    = 2629746 * second  // 30.436875 days on average
	year     = 31557600 * second // 365.2425 days on average
)

var relTimePool = sync.Pool{
	New: func() any {
		buf := make([]byte, 0, 30)
		return &buf
	},
}

func buildTime(b []byte, count int, dur string) []byte {
	b = strconv.AppendInt(b, int64(count), 10)
	b = append(b, dur...)
	if count > 1 {
		b = append(b, "s, "...)
	} else {
		b = append(b, ", "...)
	}

	return b
}

// RelTime is a wrapper around RelDuration that returns a relative
// representation of a time t.
func RelTime(t time.Time) string {
	return RelDuration(time.Since(t))
}

// RelDuration returns a relative representation of a duration d.
func RelDuration(d time.Duration) string {
	buf := relTimePool.Get().(*[]byte)
	defer relTimePool.Put(buf)
	b := (*buf)[:0]

	var parts int
	now := int(d.Abs())

	years := now / year
	if years > 0 {
		b = buildTime(b, years, " year")
		now -= years * year
		parts++
	}

	months := now / month
	if months > 0 {
		b = buildTime(b, months, " month")
		now -= months * month
		parts++
	}

	weeks := now / week
	if weeks > 0 && parts != maxParts {
		b = buildTime(b, weeks, " week")
		now -= weeks * week
		parts++
	}

	days := now / day
	if days > 0 && parts != maxParts {
		b = buildTime(b, days, " day")
		now -= days * day
		parts++
	}

	hours := now / hour
	if hours > 0 && parts != maxParts {
		b = buildTime(b, hours, " hour")
		now -= hours * hour
		parts++
	}

	minutes := now / minute
	if minutes > 0 && parts != maxParts {
		b = buildTime(b, minutes, " minute")
		now -= minutes * minute
		parts++
	}

	seconds := now / second
	if seconds > 0 && parts != maxParts {
		b = buildTime(b, seconds, " second")
	}

	if len(b) == 0 {
		return "just now"
	}

	b = append(b[:len(b)-2])
	if int(d) > 0 {
		b = append(b, " ago"...)
	}

	return *(*string)(unsafe.Pointer(&b))
}
