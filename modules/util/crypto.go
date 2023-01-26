package util

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"

	"userstyles.world/modules/config"
)

// HashIP generates a unique hash for stats.
func HashIP(key string) (string, error) {
	h := hmac.New(sha512.New, []byte(config.StatsKey))
	if _, err := h.Write([]byte(key)); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
