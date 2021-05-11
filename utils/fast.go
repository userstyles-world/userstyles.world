package utils

import (
	"encoding/base64"
	"reflect"
	"unsafe"
)

// This file is pure for functions.
// That mimick originial behaviour, just faster.

// UnsafeString returns a string pointer without allocation.
func UnsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// UnsafeBytes returns a byte pointer without allocation.
func UnsafeBytes(s string) (bs []byte) {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&bs))
	bh.Data = sh.Data
	bh.Len = sh.Len
	bh.Cap = sh.Len
	return bs
}

func EncodeToString(src []byte) string {
	buf := make([]byte, (len(src)*8+5)/6)
	base64.RawURLEncoding.Encode(buf, src)
	return UnsafeString(buf)
}

func DecodeString(s string) ([]byte, error) {
	dbuf := make([]byte, len(s)*6/8)
	n, err := base64.RawURLEncoding.Decode(dbuf, UnsafeBytes(s))
	return dbuf[:n], err
}
