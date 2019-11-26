package goacors_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/shogo82148/goacors/v2"
)

func TestEmptyOriginHeader(t *testing.T) {
	service := newService(nil)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(goacors.HeaderOrigin, "")
	rw := newTestResponseWriter()
	ctx := newContext(service, rw, req, nil)

	h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		return service.Send(ctx, http.StatusOK, "ok")
	}
	testee := goacors.New(service, &goacors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
	})(h)
	err := testee(ctx, rw, req)
	if err != nil {
		t.Error("it should not return any error but ", err)
	}
	if rw.Header().Get(goacors.HeaderAccessControlAllowOrigin) != "*" {
		t.Error("allow origin should be wild card")
	}
}

func TestOriginAllowsWildcard(t *testing.T) {
	service := newService(nil)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(goacors.HeaderOrigin, "http://someorigin.com")
	rw := newTestResponseWriter()
	ctx := newContext(service, rw, req, nil)

	h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		return service.Send(ctx, http.StatusOK, "ok")
	}
	testee := goacors.New(service, &goacors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{"*"},
	})(h)
	err := testee(ctx, rw, req)
	if err != nil {
		t.Error("it should not return any error but ", err)
	}
	if rw.Header().Get(goacors.HeaderAccessControlAllowOrigin) != req.Header.Get(goacors.HeaderOrigin) {
		t.Errorf("allow origin should be %s but %s", req.Header.Get(goacors.HeaderOrigin), rw.Header().Get(goacors.HeaderAccessControlAllowOrigin))
	}
}

func TestOrigIsNotValid(t *testing.T) {
	service := newService(nil)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(goacors.HeaderOrigin, "http://someorigin.com")
	rw := newTestResponseWriter()
	ctx := newContext(service, rw, req, nil)

	h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		return service.Send(ctx, http.StatusOK, "ok")
	}
	testee := goacors.New(service, &goacors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{"http://example.com"},
	})(h)
	err := testee(ctx, rw, req)
	if err != nil {
		t.Error("it should not return any error but ", err)
	}
	if rw.Header().Get(goacors.HeaderAccessControlAllowOrigin) != "" {
		t.Error("allow origin should be empty but ", rw.Header().Get(goacors.HeaderAccessControlAllowOrigin))
	}
}

func TestOriginAllowsFixedOrigin(t *testing.T) {
	service := newService(nil)
	fixedOrigin := "http://someorigin.com"
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(goacors.HeaderOrigin, fixedOrigin)
	rw := newTestResponseWriter()
	ctx := newContext(service, rw, req, nil)

	h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		return service.Send(ctx, http.StatusOK, "ok")
	}
	testee := goacors.New(service, &goacors.Config{
		AllowOrigins:     []string{fixedOrigin},
		ExposeHeaders:    []string{"ETag"},
		AllowCredentials: true,
	})(h)
	err := testee(ctx, rw, req)
	if err != nil {
		t.Error("it should not return any error but ", err)
	}
	if rw.Header().Get(goacors.HeaderAccessControlAllowOrigin) != fixedOrigin {
		t.Error("allow origin should be empty")
	}
	if rw.Header().Get(goacors.HeaderAccessControlExposeHeaders) != "ETag" {
		t.Error("expose header is unexpected ", rw.Header().Get(goacors.HeaderAccessControlExposeHeaders))
	}
}

func TestPreflightRequest(t *testing.T) {
	service := newService(nil)
	fixedOrigin := "http://localhost"
	req, _ := http.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set(goacors.HeaderOrigin, fixedOrigin)
	req.Header.Set(goacors.HeaderAccessControlRequestHeaders, "X-OriginalRequest")
	req.Header.Set(goacors.HeaderContentType, "application/json")
	rw := newTestResponseWriter()
	ctx := newContext(service, rw, req, nil)

	h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		return service.Send(ctx, http.StatusOK, "ok")
	}
	testee := goacors.New(service, &goacors.Config{
		AllowOrigins:     []string{fixedOrigin},
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		MaxAge:           3600,
		AllowCredentials: true,
	})(h)
	err := testee(ctx, rw, req)
	if err != nil {
		t.Error("it should not return any error but ", err)
	}
	if rw.Header().Get(goacors.HeaderAccessControlAllowOrigin) != "http://localhost" {
		t.Error("allow origin should be empty")
	}
	if rw.Header().Get(goacors.HeaderAccessControlAllowMethods) != "GET, HEAD, PUT, PATCH, POST, DELETE" {
		t.Errorf("allow method should be %q but %q", "GET, HEAD, PUT, PATCH, POST, DELETE", rw.Header().Get(goacors.HeaderAccessControlAllowMethods))
	}
	if rw.Header().Get(goacors.HeaderAccessControlAllowCredentials) != "true" {
		t.Error("allow credentials should be true")
	}
	if rw.Header().Get(goacors.HeaderAccessControlMaxAge) != "3600" {
		t.Error("access control max age should be 3600 but ", rw.Header().Get(goacors.HeaderAccessControlMaxAge))
	}
	if rw.Header().Get(goacors.HeaderAccessControlAllowHeaders) != "X-OriginalRequest" {
		t.Error("access control allow headers should be 'X-OriginalRequest' but ", rw.Header().Get(goacors.HeaderAccessControlAllowHeaders))
	}

	// StatusNoContent does not allow body
	if rw.Status != http.StatusNoContent {
		t.Errorf("the status should be %d, got %d", http.StatusNoContent, rw.Status)
	}
	if len(rw.Body) != 0 {
		t.Errorf("the length of the body should be 0, got %d", len(rw.Body))
	}
}

func TestNotGivenAllowHeaderOnRequest(t *testing.T) {
	service := newService(nil)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(goacors.HeaderOrigin, "http://localhost")
	rw := newTestResponseWriter()
	ctx := newContext(service, rw, req, nil)

	h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		return service.Send(ctx, http.StatusOK, "ok")
	}
	testee := goacors.New(service, &goacors.Config{
		AllowCredentials: true,
		AllowOrigins:     []string{"http://example.com"},
	})(h)
	err := testee(ctx, rw, req)
	if err != nil {
		t.Fatal("it should not return any error but ", err)
	}
	if rw.Header().Get(goacors.HeaderAccessControlAllowOrigin) != "" {
		t.Error("allow origin should be empty")
	}
}

func TestExecuteWithSkipper(t *testing.T) {
	service := newService(nil)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(goacors.HeaderOrigin, "mismatchedhost")
	rw := newTestResponseWriter()
	ctx := newContext(service, rw, req, nil)

	h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		return service.Send(ctx, http.StatusOK, "ok")
	}
	testee := goacors.New(service, &goacors.Config{
		Skipper: func(c context.Context, rw http.ResponseWriter, req *http.Request) bool {
			return true
		},
		AllowCredentials: true,
		AllowOrigins:     []string{"http://example.com"},
	})(h)
	err := testee(ctx, rw, req)
	if err != nil {
		t.Fatal("it should not return any error but ", err)
	}
	if rw.Header().Get(goacors.HeaderAccessControlAllowOrigin) != "" {
		t.Error("allow origin should be empty")
	}
}

func TestRequestGetWithOrigin(t *testing.T) {
	service := newService(nil)
	fixedOrigin := "http://localhost"
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(goacors.HeaderOrigin, fixedOrigin)
	req.Header.Set(goacors.HeaderContentType, "application/json")
	rw := newTestResponseWriter()
	ctx := newContext(service, rw, req, nil)

	h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		return service.Send(ctx, http.StatusOK, "ok")
	}
	testee := goacors.New(service, &goacors.Config{
		AllowOrigins:     []string{fixedOrigin},
		AllowCredentials: true,
	})(h)
	err := testee(ctx, rw, req)
	if err != nil {
		t.Error("it should not return any error but ", err)
	}
	if rw.Header().Get(goacors.HeaderAccessControlAllowOrigin) != "http://localhost" {
		t.Error("allow origin should be empty")
	}
}

func TestAddedAllowOrigHeader(t *testing.T) {
	service := newService(nil)
	req, _ := http.NewRequest(http.MethodOptions, "/", nil)
	req.Header.Set(goacors.HeaderOrigin, "http://someorigin.com")
	rw := newTestResponseWriter()
	ctx := newContext(service, rw, req, nil)

	h := func(ctx context.Context, rw http.ResponseWriter, req *http.Request) error {
		return service.Send(ctx, http.StatusOK, "ok")
	}
	testee := goacors.New(service, &goacors.Config{
		AllowCredentials: true,
		AllowHeaders:     []string{"X-OrigHeader"},
	})(h)
	err := testee(ctx, rw, req)
	if err != nil {
		t.Error("it should not return any error but ", err)
	}
	if rw.Header().Get(goacors.HeaderAccessControlAllowHeaders) != "X-OrigHeader" {
		t.Error("allow origin should be empty")
	}

	// StatusNoContent does not allow body
	if rw.Status != http.StatusNoContent {
		t.Errorf("the status should be %d, got %d", http.StatusNoContent, rw.Status)
	}
	if len(rw.Body) != 0 {
		t.Errorf("the length of the body should be 0, got %d", len(rw.Body))
		t.Log(string(rw.Body))
	}
}
