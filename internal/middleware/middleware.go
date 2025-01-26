package middleware

import "net/http"

// RoundTripFunc is how you add middleware (interceptor) to the client.
type RoundTripFunc func(req *http.Request) (*http.Response, error)

func (r RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return r(req)
}

// Applies middlewares, they will be called in the order
func ApplyMiddlewares(c *http.Client, middlewares ...func(http.RoundTripper) http.RoundTripper) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		c.Transport = middlewares[i](c.Transport)
	}
}
