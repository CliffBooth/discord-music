package middleware

import (
	"fmt"
	"net/http"
)

// TODO: use mutex when updating cache here
type R struct {
}

func RateLimitMiddleware(next http.RoundTripper) http.RoundTripper {
	return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		fmt.Println("ratelimit: before calling next()")
		resp, err := next.RoundTrip(req)
		fmt.Println("ratelimit: after calling next()")
		return resp, err
	})
}
