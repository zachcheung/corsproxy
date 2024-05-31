package corsproxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/rs/cors"
)

// Options defines the configuration for the CORS proxy.
type Options struct {
	cors.Options
	// AllowedTargets is a list of targets a cross-domain request can reach.
	// If the special "*" value is present in the list, all targets will be allowed except private network targets if `AllowPrivateNetworkTarget` is false.
	// A target may contain a wildcard (*) to replace 0 or more characters (i.e.: `http://*.domain.com`).
	// Only one wildcard can be used per target.
	// Default value is ["*"]
	AllowedTargets []string
	// AllowPrivateNetworkTarget indicates whether to accept private network targets.
	// If a private network target has been added to `AllowedTargets`, there is no need to set `AllowPrivateNetworkTarget` explicitly.
	AllowPrivateNetworkTarget bool
}

// CorsProxy holds the configuration and state for the CORS proxy.
type CorsProxy struct {
	cors                      *cors.Cors
	allowedTargets            []string
	allowedWTargets           []wildcard
	allowedTargetsAll         bool
	allowPrivateNetworkTarget bool
}

// New creates a new CorsProxy instance with the provided options.
func New(options Options) *CorsProxy {
	cp := &CorsProxy{
		allowPrivateNetworkTarget: options.AllowPrivateNetworkTarget,
	}

	// Handle allowed targets
	switch {
	case len(options.AllowedTargets) == 0:
		// Default is to allow all targets except private network targets if AllowPrivateNetworkTarget is false
		cp.allowedTargetsAll = true
	default:
		for _, target := range options.AllowedTargets {
			target = strings.ToLower(target)
			if target == "*" {
				// If "*" is present in the list, allow all targets except private network targets if AllowPrivateNetworkTarget is false
				cp.allowedTargetsAll = true
				cp.allowedTargets = nil
				cp.allowedWTargets = nil
				break
			} else if i := strings.IndexByte(target, '*'); i >= 0 {
				// Split the target into prefix and suffix based on the wildcard '*'
				w := wildcard{target[:i], target[i+1:]}
				cp.allowedWTargets = append(cp.allowedWTargets, w)
			} else if target != "" {
				cp.allowedTargets = append(cp.allowedTargets, target)
			}

			if !cp.allowPrivateNetworkTarget {
				remote, err := url.Parse(target)
				if err != nil {
					panic(fmt.Sprintf("Invalid target %s in AllowedTargets: %v", target, err))
				}
				private, parsed := isPrivateAddr(remote.Host)
				if parsed && private {
					cp.allowPrivateNetworkTarget = true
				}
			}
		}
	}

	cp.cors = cors.New(options.Options)
	return cp
}

// Handler returns an HTTP handler that proxies requests with CORS support.
func (cp *CorsProxy) Handler() http.Handler {
	return cp.cors.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := strings.TrimPrefix(r.URL.Path, "/")
		if target == "" || target == "favicon.ico" {
			w.WriteHeader(http.StatusOK)
			return
		}

		targetStr := strings.TrimPrefix(r.URL.String(), "/")

		// Check if the target URL is allowed
		allowed, statusCode := cp.isTargetAllowed(r, target)
		if !allowed {
			var msg string
			switch statusCode {
			case http.StatusBadRequest:
				msg = "Invalid URL"
			case http.StatusForbidden:
				msg = "Forbidden"
			default:
				msg = "Unknown error"
			}
			http.Error(w, msg, statusCode)
			return
		}

		// Parse and normalize the target URL
		remote, err := NormalizeParseURL(targetStr)
		if err != nil {
			http.Error(w, "Invalid target URL", http.StatusBadRequest)
			return
		}

		// Create the reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = r.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = remote.Path
			req.URL.RawQuery = remote.RawQuery
		}

		// Serve the request using the reverse proxy
		proxy.ServeHTTP(w, r)
	}))
}

// isTargetAllowed checks if the target URL is allowed based on the proxy's configuration.
func (cp *CorsProxy) isTargetAllowed(r *http.Request, target string) (allowed bool, statusCode int) {
	// Parse and normalize the target URL
	remote, err := NormalizeParseURL(target)
	if err != nil {
		return false, http.StatusBadRequest
	}

	// Ensure the target URL has a valid scheme
	if remote.Scheme != "https" && remote.Scheme != "http" {
		return false, http.StatusBadRequest
	}

	if !cp.allowPrivateNetworkTarget {
		private, parsed := isPrivateAddr(remote.Host)
		if parsed && private {
			return false, http.StatusForbidden
		}
	}

	// If all targets are allowed
	if cp.allowedTargetsAll {
		return true, http.StatusOK
	}

	// Normalize the target URL for comparison
	normalizedTarget := strings.TrimSuffix(target, "/")
	// Check against allowed targets
	for _, t := range cp.allowedTargets {
		normalizedT := strings.TrimSuffix(t, "/")
		if normalizedTarget == normalizedT || strings.HasPrefix(target, normalizedT+"/") {
			return true, http.StatusOK
		}
	}

	// Check against wildcard allowed targets
	for _, w := range cp.allowedWTargets {
		if w.match(target) {
			return true, http.StatusOK
		}
	}

	return false, http.StatusForbidden
}
