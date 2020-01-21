package goacors

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/shogo82148/goa-v1"
)

// New creates middleware with configure for this
func New(service *goa.Service, conf *Config) goa.Middleware {
	// validate allowed origin configure
	allowAnyOrigin := false
	allowOrigins := make([]originType, len(conf.AllowOrigins))
	for i, origin := range conf.AllowOrigins {
		if origin == "*" {
			allowAnyOrigin = true
			break
		}
		o, err := parseOrigin(origin)
		if err != nil {
			panic("invalid allowed origin: " + origin)
		}
		allowOrigins[i] = o
	}

	skipper := conf.Skipper
	allowMethods := strings.Join(conf.AllowMethods, ", ")
	allowHeaders := strings.Join(conf.AllowHeaders, ", ")
	exposeHeaders := strings.Join(conf.ExposeHeaders, ", ")
	allowCredentials := conf.AllowCredentials
	var maxAge string
	if conf.MaxAge > 0 {
		maxAge = strconv.Itoa(conf.MaxAge)
	}

	return func(next goa.Handler) goa.Handler {
		return func(c context.Context, rw http.ResponseWriter, req *http.Request) error {
			// Skipper
			if skipper != nil && skipper(c, rw, req) {
				return next(c, rw, req)
			}

			h := rw.Header()

			// Check the origin of the request is allowed
			var allowedOrigin string
			if allowAnyOrigin {
				if allowCredentials {
					// https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS
					// When responding to a credentialed request, the server must specify an origin in the value of
					// the Access-Control-Allow-Origin header, instead of specifying the "*" wildcard.
					allowedOrigin = req.Header.Get(HeaderOrigin)
				} else {
					allowedOrigin = "*"
				}
			} else {
				origin := req.Header.Get(HeaderOrigin)
				if allowed(origin, allowOrigins, allowCredentials) {
					allowedOrigin = origin
				}
			}

			if req.Method != http.MethodOptions {
				// handle normal requests
				h.Add(HeaderVary, HeaderOrigin)
				if allowedOrigin != "" {
					h.Set(HeaderAccessControlAllowOrigin, allowedOrigin)
				}
				if allowCredentials {
					h.Set(HeaderAccessControlAllowCredentials, "true")
				}
				if exposeHeaders != "" {
					h.Set(HeaderAccessControlExposeHeaders, exposeHeaders)
				}
				return next(c, rw, req)
			}

			// handle preflight requests
			h.Add(HeaderVary, HeaderOrigin)
			h.Add(HeaderVary, HeaderAccessControlRequestMethod)
			h.Add(HeaderVary, HeaderAccessControlRequestHeaders)
			h.Set(HeaderAccessControlAllowOrigin, allowedOrigin)
			h.Set(HeaderAccessControlAllowMethods, allowMethods)
			if allowCredentials {
				h.Set(HeaderAccessControlAllowCredentials, "true")
			}
			if allowHeaders != "" {
				h.Set(HeaderAccessControlAllowHeaders, allowHeaders)
			} else {
				header := req.Header.Get(HeaderAccessControlRequestHeaders)
				if header != "" {
					h.Set(HeaderAccessControlAllowHeaders, header)
				}
			}

			if maxAge != "" {
				h.Set(HeaderAccessControlMaxAge, maxAge)
			}
			rw.WriteHeader(http.StatusNoContent)
			return nil
		}
	}

}
