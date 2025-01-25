package discord

import (
	httpClient "discord-music/internal/httpClient"
	"fmt"
	"net/http"
)

type R struct {
}

func RateLimitMiddleware(next httpClient.Middleware) httpClient.Middleware {
	return func(req *http.Request) (*http.Response, error) {
		fmt.Println("ratelimit: before calling next()")
		resp, err := next(req)
		fmt.Println("ratelimit: after calling next()")
		return resp, err
	}
}
