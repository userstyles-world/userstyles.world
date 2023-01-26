package util

import (
	"testing"
)

func TestHashIP(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		id       string
		ip       string
		expected string
	}{
		{"TestNormalCase", "test", "test", "4942c2726e0c6ad350c2ecebc1ef37e74d6802567c7b51df0c08968a00f4dd8df541baf0495eb2127a53aca81755f3ec3e3434db49ad326dd7b607809a564acf"},
		{"TestRealIP", "Secrecy", "192.168.1.1", "cb18a873a9195e00320830f11739227df9b4d6a9462c58a1b3298a17297dc431b4880c57375d22a58a6149e3586a90e63075603c02c55db6ae3459039312ae3e"},
		{"TestEmptyIP", "Secrecy", "", "12f996268829112ef8f87e7b6e0dc65b8e64642f13a134f49740c7769062ecc259cc740c7cb7e7f827e202d347dbf4dd88d55b5d495e1bf7006b808c692a68bd"},
		{"TestEmptyID", "", "192.168.1.1", "1778721963e12865be624b66e70103c14cbb1c2ccc494a990927edb89f67eb58232ac8bad231c75192643808b75a6d86228b5b03da2743af748c0cc1be68aa9f"},
	}

	for _, c := range cases {
		actual, err := CreateHashedRecord(c.ip + " " + c.id)
		if err != nil {
			t.Error(err)
		}
		if actual != c.expected {
			t.Fatalf("%s: expected: %s got: %s",
				c.name, c.expected, actual)
		}
	}
}
