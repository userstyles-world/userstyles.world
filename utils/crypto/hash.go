package crypto

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"

	"userstyles.world/modules/config"
)

// CreateHashedRecord generates a unique hash for stats.
func CreateHashedRecord(key string) (string, error) {
	// Generate unique hash.
	h := hmac.New(sha512.New, []byte(config.StatsKey))
	if _, err := h.Write([]byte(key)); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
