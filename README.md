[![Build Status](https://github.com/shogo82148/goacors-v1/workflows/Test/badge.svg?branch=master)](https://github.com/shogo82148/goacors-v1/actions)
[![Coverage Status](https://coveralls.io/repos/github/shogo82148/goacors-v1/badge.svg?branch=master&service=github)](https://coveralls.io/github/shogo82148/goacors-v1?branch=master) [![PkgGoDev](https://pkg.go.dev/badge/github.com/shogo82148/goacors-v1?tab=doc)](https://pkg.go.dev/badge/github.com/shogo82148/goacors-v1?tab=doc)

# goacors-v1

a cors-header middleware for goa(https://github.com/shogo82148/goa-v1).
This is a fork of https://github.com/istyle-inc/goacors

# how to use

1. `go get github.com/shogo82148/goacors-v1`
2. write your main.go generated automatically from goagen.

```go
service.Use(goacors.New(service, &goacors.Config{
	AllowOrigins: []string{"http://example.com"},
	AllowMethods: []string{http.MethodGet},
}))
```
