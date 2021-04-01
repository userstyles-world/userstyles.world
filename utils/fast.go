package utils

import (
	"encoding/base64"
	"reflect"
	"unsafe"
)

// This file is pure for functions.
// That mimick originial behaviour, just faster.

func B2s(msg []byte) string {
	return *(*string)(unsafe.Pointer(&msg))
}

// Note it may break if string and/or slice header will change
// in the future go versions.
func S2b(str string) []byte {
	var b []byte
	byteHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	byteHeader.Data = (*reflect.StringHeader)(unsafe.Pointer(&str)).Data

	// This reference is important as without it their is an chance
	// That the str get GC'ed.
	l := len(str)
	byteHeader.Len = l
	byteHeader.Cap = l

	return b
}

func EncodeToString(src []byte) string {
	buf := make([]byte, (len(src)*8+5)/6)
	base64.RawURLEncoding.Encode(buf, src)
	return B2s(buf)
}

func DecodeString(s string) ([]byte, error) {
	dbuf := make([]byte, len(s)*6/8)
	n, err := base64.RawURLEncoding.Decode(dbuf, S2b(s))
	return dbuf[:n], err
}
