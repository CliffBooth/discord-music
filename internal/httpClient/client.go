package httpclient

import "net/http"

// this is the wrapping of http.Client but with middleware

/*
middleware:
	type handler func(request) response

	func(f handler) handler

	func myMiddleware(next handler) handler {
		return func(request) response {
			//do something...
			response = next(requst)
			//do something...
			reuturn response
		}
	}
*/

type Client struct {
	c  http.Client
	do func(req *http.Request) (*http.Response, error)
}

type Middleware func(req *http.Request) (*http.Response, error)

// Applies middlewares, they will be called in the order
func (c *Client) ApplyMiddlewares(middlewares ...func(m Middleware) Middleware) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		c.do = middlewares[i](c.do)
	}
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	return c.do(req)
}

func New(c http.Client) *Client {
	return &Client{
		c:  c,
		do: c.Do,
	}
}
