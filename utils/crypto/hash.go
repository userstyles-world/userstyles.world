package crypto

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"

	"userstyles.world/config"
)

func CreateHashedRecord(id, ip string) (string, error) {
	// Merge it here before using, otherwise errors occur.
	record := ip + " " + id

	// Generate unique hash.
	h := hmac.New(sha512.New, []byte(config.STATS_KEY))
	if _, err := h.Write([]byte(record)); err != nil {
		return "", err
	}
	s := hex.EncodeToString(h.Sum(nil))

	return s, nil
}
