package goacors

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// https://developer.mozilla.org/en-US/docs/Glossary/Origin
// > Web content's origin is defined by the scheme (protocol), host (domain), and port of the URL used to access it.
type originType struct {
	scheme string
	host   string
	port   int
}

func parseOrigin(s string) (originType, error) {
	var origin originType
	u, err := url.Parse(s)
	if err != nil {
		return originType{}, err
	}

	switch u.Scheme {
	case "http":
		origin.scheme = "http"
		origin.port = 80 // default port for http
	case "https":
		origin.scheme = "https"
		origin.port = 443 // default port for https
	case "":
		return originType{}, fmt.Errorf("goacors: scheme is required: %s", s)
	default:
		// only support http and https
		return originType{}, fmt.Errorf("goacors: unknown scheme: %s", u.Scheme)
	}

	// host is case insensitive
	origin.host = strings.ToLower(u.Hostname())

	if port := u.Port(); port != "" {
		num, err := strconv.Atoi(port)
		if err != nil {
			return originType{}, err
		}
		origin.port = num
	}

	return origin, nil
}

func match(origin, allowed originType) bool {
	if origin.scheme != allowed.scheme {
		return false
	}
	if origin.port != allowed.port {
		return false
	}

	// handle wildcard domain
	for strings.HasPrefix(allowed.host, "*.") {
		idx := strings.Index(origin.host, ".")
		if idx <= 0 {
			return false
		}
		origin.host = origin.host[idx+1:]
		allowed.host = allowed.host[len("*."):]
	}
	return origin.host == allowed.host
}

func allowed(origin string, allowedOrigins []originType, allowCredentials bool) bool {
	o, err := parseOrigin(origin)
	if err != nil {
		return false
	}
	for _, allowedOrigin := range allowedOrigins {
		if match(o, allowedOrigin) {
			return true
		}
	}
	return false
}
