package corsproxy

import (
	"net/netip"
	"net/url"
	"strings"
)

// wildcard represents a pattern for wildcard matching in allowed targets.
type wildcard struct {
	prefix string
	suffix string
}

// match checks if the target URL matches the wildcard pattern.
func (w wildcard) match(s string) bool {
	if len(s) >= len(w.prefix)+len(w.suffix) {
		if after, found := strings.CutPrefix(s, w.prefix); found {
			return strings.HasSuffix(after, w.suffix) || strings.Contains(after, w.suffix+"/")
		}
	}

	return false
}

// isPrivateAddr checks if the given host is a private or loopback IP address.
func isPrivateAddr(host string) (private bool, parsed bool) {
	var (
		addr     netip.Addr
		addrPort netip.AddrPort
		err      error
	)

	if strings.Contains(host, ":") {
		// Parse as address:port
		addrPort, err = netip.ParseAddrPort(host)
		if err == nil {
			addr = addrPort.Addr()
		}
	} else {
		// Parse as address only
		addr, err = netip.ParseAddr(host)
	}

	if err != nil {
		return false, false
	}

	// Check if the address is private or loopback
	return addr.IsPrivate() || addr.IsLoopback(), true
}

// stripURLQuery takes a raw URL string, normalizes it, and returns the URL without the query string.
func stripURLQuery(rawURL string) (string, error) {
	normalizedURL, err := normalizeURL(rawURL)
	if err != nil {
		return "", err
	}

	return strings.SplitN(normalizedURL, "?", 2)[0], nil
}

// normalizeURL takes a raw URL string and normalizes it by converting the host to lowercase.
func normalizeURL(rawURL string) (string, error) {
	parsedURL, err := normalizeParseURL(rawURL)
	if err != nil {
		return "", err
	}

	return parsedURL.String(), nil
}

// normalizeParseURL takes a raw URL string, parses it into a *url.URL, and normalizes it by converting the host to lowercase.
func normalizeParseURL(rawURL string) (*url.URL, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	parsedURL.Host = strings.ToLower(parsedURL.Host)

	return parsedURL, nil
}
