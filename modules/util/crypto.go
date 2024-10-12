package util

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"sync"
	"unsafe"

	"userstyles.world/modules/config"
)

var (
	bufHashPool = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 128)
			return &b
		},
	}
	hmacPool sync.Pool
)

func InitPool(secrets *config.SecretsConfig) {
	hmacPool = sync.Pool{
		New: func() interface{} {
			return hmac.New(sha512.New, []byte(secrets.StatsKey))
		},
	}
}

// HashIP generates a unique hash for stats.
func HashIP(key string) (string, error) {
	bp := bufHashPool.Get().(*[]byte)
	defer bufHashPool.Put(bp)
	b := *bp

	h := hmacPool.Get().(hash.Hash)
	defer hmacPool.Put(h)
	h.Reset()

	if _, err := h.Write([]byte(key)); err != nil {
		return "", err
	}

	hex.Encode(b, h.Sum(nil))

	return *(*string)(unsafe.Pointer(&b)), nil
}
