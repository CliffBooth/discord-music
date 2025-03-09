package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"discord-music/internal/log"
)

func GetLogMiddleware(logger log.Logger) func(http.RoundTripper) http.RoundTripper {
	return func(next http.RoundTripper) http.RoundTripper {
		return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			printRequest(logger, req)
			resp, err := next.RoundTrip(req)
			if err == nil {
				printResponse(logger, resp)
			} else {
				logger.Debugf("[log middleware] error making RoundTrip: %v\n", err)
			}
			return resp, err
		})
	}
}

func headersAsString(h http.Header) string {
	headers := make([]string, 0, len(h))
	for k := range h {
		headers = append(headers, k)
	}
	sort.Strings(headers)

	sb := strings.Builder{}

	for _, k := range headers {
		sb.WriteString(fmt.Sprintf("%s: %s; ", k, h[k]))
	}

	return sb.String()
}

func printRequest(l log.Logger, req *http.Request) {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s %v ", req.Method, req.URL))
	sb.WriteString(headersAsString(req.Header))

	l.Debug("[log middleware] reqeust: ", sb.String())
}

func printResponse(l log.Logger, resp *http.Response) {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%v ", resp.StatusCode))

	sb.WriteString(headersAsString(resp.Header))

	newBody := bytes.NewBuffer([]byte{})
	r := io.TeeReader(resp.Body, newBody)
	body, err := io.ReadAll(r) // when we do this, the body will be written to newBody

	resp.Body = io.NopCloser(newBody)

	if err != nil {
		l.Debug("[log middleware] error reading body")
		return
	}
	sb.WriteString(string(body))

	l.Debug("[log middleware] response: ", sb.String())
}
