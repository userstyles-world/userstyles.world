package utils

import (
	"encoding/base64"
	"reflect"
	"unsafe"

	"github.com/ohler55/ojg/oj"
)

// This file is functions that mimic original behaviour, just faster.

var jsonEncoderOptions = &oj.Options{OmitNil: true, UseTags: true}

// UnsafeString returns a string pointer without allocation.
func UnsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// UnsafeClone copies the current string to a new string with 1 allocation.
// https://go-review.googlesource.com/c/go/+/345849/1/src/strings/clone.go
func UnsafeClone(s string) string {
	b := make([]byte, len(s))
	copy(b, s)
	return *(*string)(unsafe.Pointer(&b))
}

// UnsafeBytes returns a byte pointer without allocation.
func UnsafeBytes(s string) (bs []byte) {
	return unsafe.Slice((*byte)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data)), len(s))
}

func EncodeToString(src []byte) string {
	buf := make([]byte, (len(src)*8+5)/6)
	base64.RawURLEncoding.Encode(buf, src)
	return UnsafeString(buf)
}

func decodeBase64(s string) ([]byte, error) {
	dbuf := make([]byte, len(s)*6/8)
	n, err := base64.RawURLEncoding.Decode(dbuf, UnsafeBytes(s))
	return dbuf[:n], err
}

func JSONEncoder(value interface{}) ([]byte, error) {
	return oj.Marshal(value, jsonEncoderOptions)
}
