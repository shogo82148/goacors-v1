package goacors

import (
	"net/http"
)

const (
	// HeaderVary "Vary"
	HeaderVary = "Vary"
	// HeaderOrigin "Origin"
	HeaderOrigin = "Origin"
	// HeaderAccessControlRequestMethod "Access-Control-Request-Method"
	HeaderAccessControlRequestMethod = "Access-Control-Request-Method"
	// HeaderAccessControlRequestHeaders "Access-Control-Request-Headers"
	HeaderAccessControlRequestHeaders = "Access-Control-Request-Headers"
	// HeaderAccessControlAllowOrigin Access-Control-Allow-Origin"
	HeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"
	// HeaderAccessControlAllowMethods "Access-Control-Allow-Methods"
	HeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"
	// HeaderAccessControlAllowHeaders "Access-Control-Allow-Headers"
	HeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"
	// HeaderAccessControlAllowCredentials "Access-Control-Allow-Credentials"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	// HeaderAccessControlExposeHeaders "Access-Control-Expose-Headers"
	HeaderAccessControlExposeHeaders = "Access-Control-Expose-Headers"
	// HeaderAccessControlMaxAge "Access-Control-Max-Age"
	HeaderAccessControlMaxAge = "Access-Control-Max-Age"
	// HeaderContentType "Content-Type"
	HeaderContentType = "Content-Type"
)

// DomainStrategy defined identify how handle (judge match with origin or not) domain
type DomainStrategy int

const (
	// AllowStrict strict mode (completely same origin or wild card or null)
	AllowStrict DomainStrategy = iota
	// AllowIntermediateMatch intermediate-match (such as subdomain like '*.example.com')
	AllowIntermediateMatch
)

// Skipper defines a function to skip middleware. Returning true skips processing
// the middleware.
type Skipper func(c context.Context, rw http.ResponseWriter, req *http.Request) bool

// Config is a config for the CORS middleware.
type Config struct {
	// Skipper defines a function to skip middleware.
	Skipper Skipper

	DomainStrategy DomainStrategy

	// AllowOrigin defines a list of origins that may access the resource.
	// Default value is an empty list, any origin can not access.
	AllowOrigins []string

	// AllowMethods defines a list methods allowed when accessing the resource.
	// This is used in response to a preflight request.
	// Default value is an empty list, any method is not allowed.
	AllowMethods []string

	// AllowHeaders defines a list of request headers that can be used when
	// making the actual request. This in response to a preflight request.
	AllowHeaders []string

	// AllowCredentials indicates whether or not the response to the request
	// can be exposed when the credentials flag is true. When used as part of
	// a response to a preflight request, this indicates whether or not the
	// actual request can be made using credentials.
	// Default value is false.
	AllowCredentials bool

	// ExposeHeaders defines a whitelist headers that clients are allowed to
	// access.
	// Default value is an empty list.
	ExposeHeaders []string

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached.
	// The default value is 0, the preflight request can not be cached.
	MaxAge int
}
