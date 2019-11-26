package goacors

import (
	"reflect"
	"testing"
)

func TestParseOrigin(t *testing.T) {
	testcases := []struct {
		in  string
		out originType
		err bool
	}{
		{
			in: "http://example.com",
			out: originType{
				scheme: "http",
				host:   "example.com",
				port:   80,
			},
		},
		{
			in: "https://example.com",
			out: originType{
				scheme: "https",
				host:   "example.com",
				port:   443,
			},
		},

		// origin is case insensitive
		{
			in: "HTTP://EXAMPLE.COM",
			out: originType{
				scheme: "http",
				host:   "example.com",
				port:   80,
			},
		},

		// set port number
		{
			in: "http://example.com:8080",
			out: originType{
				scheme: "http",
				host:   "example.com",
				port:   8080,
			},
		},

		{
			in:  "example.com",
			err: true,
		},
		{
			in:  "ftp://example.com",
			err: true,
		},
	}

	for i, tc := range testcases {
		origin, err := parseOrigin(tc.in)
		if err != nil {
			if !tc.err {
				t.Errorf("%d: want not error, got error: %v", i, err)
			}
		} else {
			if tc.err {
				t.Errorf("%d: want error, got not error", i)
			} else if !reflect.DeepEqual(origin, tc.out) {
				t.Errorf("%d: want %+v, got %+v", i, tc.out, origin)
			}
		}
	}
}

func TestMatch(t *testing.T) {
	testcases := []struct {
		origin  string
		allowed string
		want    bool
	}{
		{
			origin:  "http://example.com",
			allowed: "http://example.com",
			want:    true,
		},
		{
			origin:  "http://example.com",
			allowed: "http://example.com:80",
			want:    true,
		},
		{
			origin:  "http://example.com",
			allowed: "https://example.com",
			want:    false,
		},

		// wildcard domains
		{
			origin:  "http://foo.example.com",
			allowed: "http://*.example.com",
			want:    true,
		},
		{
			origin:  "http://example.com",
			allowed: "http://*.example.com",
			want:    false,
		},
		{
			origin:  "http://foo.bar.example.com",
			allowed: "http://*.example.com",
			want:    false,
		},
		{
			origin:  "http://foo.bar.example.com",
			allowed: "http://*.*.example.com",
			want:    true,
		},
	}

	for i, tc := range testcases {
		origin, err := parseOrigin(tc.origin)
		if err != nil {
			t.Errorf("%d: error %v", i, err)
			continue
		}
		allowed, err := parseOrigin(tc.allowed)
		if err != nil {
			t.Errorf("%d: error %v", i, err)
			continue
		}
		got := match(origin, allowed)
		if got != tc.want {
			t.Errorf("%d: want %v, got %v", i, tc.want, got)
		}
	}
}
