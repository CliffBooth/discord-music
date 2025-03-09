package middleware

import (
	"net/http"
)

// TODO: use mutex when updating cache here
type R struct {
}

func RateLimitMiddleware(next http.RoundTripper) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		resp, err := next.RoundTrip(req)
		return resp, err
	})
}
