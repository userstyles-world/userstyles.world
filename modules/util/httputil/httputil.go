// Package httputil provides helper functions for HTTP facilities.
package httputil

// ProxyHeader will return proper IP depending on the environment.
func ProxyHeader(production bool) string {
	if production {
		return "X-Real-IP"
	}

	return ""
}
