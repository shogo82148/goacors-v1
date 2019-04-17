package goacors

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/goadesign/goa"
)

// New return middleware implements checking cors with default config
func New(service *goa.Service) goa.Middleware {
	return WithConfig(service, nil)
}

// WithConfig create middleware with configure for this
func WithConfig(service *goa.Service, conf *Config) goa.Middleware {
	if conf == nil {
		conf = DefaultConfig
	}

	skipper := conf.Skipper
	if len(conf.AllowOrigins) == 0 {
		conf.AllowOrigins = DefaultConfig.AllowOrigins
	}
	allowMethods := strings.Join(conf.AllowMethods, ", ")
	allowHeaders := strings.Join(conf.AllowHeaders, ", ")
	exposeHeaders := strings.Join(conf.ExposeHeaders, ", ")
	var maxAge string
	if conf.MaxAge > 0 {
		maxAge = strconv.Itoa(conf.MaxAge)
	}

	var om OriginMatcher
	switch conf.DomainStrategy {
	case AllowIntermediateMatch:
		om = newInterMediateMatcher(conf)
	case AllowStrict:
		om = newStrictOriginMatcher(conf)
	default:
		panic(fmt.Errorf("goacors: invalid domain strategy: %d", conf.DomainStrategy))
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
			allowedOrigin, _ := om.FindMatchedOrigin(conf.AllowOrigins, origin)

			if req.Method != http.MethodOptions {
				// handle normal requests
				h.Add(HeaderVary, HeaderOrigin)
				h.Set(HeaderAccessControlAllowOrigin, allowedOrigin)
				if conf.AllowCredentials && allowedOrigin != "*" && allowedOrigin != "" {
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
			if conf.AllowCredentials && allowedOrigin != "*" && allowedOrigin != "" {
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
