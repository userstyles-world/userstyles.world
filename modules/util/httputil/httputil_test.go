package httputil

import "testing"

func TestProxyHeader(t *testing.T) {
	t.Parallel()

	development := false
	if ProxyHeader(development) != "" {
		t.Fatal("should return an empty string")
	}

	production := true
	if ProxyHeader(production) != "X-Real-IP" {
		t.Fatal("should return X-Real-IP")
	}
}
