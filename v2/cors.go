package goacors

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/goadesign/goa"
)

// New creates middleware with configure for this
func New(service *goa.Service, conf *Config) goa.Middleware {
	skipper := conf.Skipper
	allowOrigins := make([]string, len(conf.AllowOrigins))
	copy(allowOrigins, conf.AllowOrigins)
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

			// Check allowed origins
			origin := req.Header.Get(HeaderOrigin)
			allowedOrigin, _ := findMatchedOrigin(allowOrigins, origin, allowCredentials)

			if req.Method != http.MethodOptions {
				// handle normal requests
				h.Add(HeaderVary, HeaderOrigin)
				h.Set(HeaderAccessControlAllowOrigin, allowedOrigin)
				if allowCredentials && allowedOrigin != "*" && allowedOrigin != "" {
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
			if allowCredentials && allowedOrigin != "*" && allowedOrigin != "" {
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
