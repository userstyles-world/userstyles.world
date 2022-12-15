package utils

import "net"

// IsLoopback returns if the address resolves to a local address.
func IsLoopback(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}
	if host == "localhost" || host == "127.0.0.1" || host == "::1" {
		return true
	}

	// Check for loopback network.
	ips, err := net.LookupIP(host)
	if err != nil {
		return false
	}
	for _, ip := range ips {
		if !ip.IsLoopback() {
			return false
		}
	}
	return true
}

// IsLocal returns if the address resolves to a local address in production.
func IsLocal(production bool, addr string) bool {
	// Because IP addresses are proxied in production, Fiber is returning an
	// empty string if "X-Real-IP" header is absent.  Since our main use-case
	// for this is to run `go tool pprof`, we can ignore this check altogether.
	if production && addr == "" {
		return true
	}
	return IsLoopback(addr)
}
