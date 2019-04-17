package goacors

func findMatchedOrigin(allowedOrigins []string, origin string, allowCredentials bool) (foundOne string, found bool) {
	for _, o := range allowedOrigins {
		if foundOne, found = innerMatcher(o, origin, allowCredentials); found {
			return
		}
	}
	return
}

func innerMatcher(allowedOrigin string, origin string, allowCredentials bool) (filteredOrigin string, ok bool) {
	if allowedOrigin == "*" && allowCredentials && origin != "" {
		return origin, true
	}
	if allowedOrigin == "*" || allowedOrigin == origin {
		return allowedOrigin, true
	}
	return
}
