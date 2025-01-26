package middleware

import (
	"fmt"
	"io"
	"net/http"
	"sort"
)

// TODO: add your logger here, and print with logLevel = debug
func LogMiddleware(next http.RoundTripper) http.RoundTripper {

	return RoundTripFunc(func(req *http.Request) (*http.Response, error) {
		resp, err := next.RoundTrip(req)
		printRequest(req)
		if err == nil {
			printResponse(resp)
		} else {
			fmt.Printf("error occured: %v\n", err)
		}
		return resp, err
	})
}

func printRequest(req *http.Request) {
	fmt.Printf("%s %s\n", req.Method, req.URL)
	fmt.Println()
	headers := make([]string, 0, len(req.Header))
	for k := range req.Header {
		headers = append(headers, k)
	}
	sort.Strings(headers)
	for _, k := range headers {
		fmt.Printf("%s: %v\n", k, req.Header[k])
	}
}

func printResponse(resp *http.Response) {
	fmt.Println(resp.StatusCode)
	fmt.Println()
	headers := make([]string, 0, len(resp.Header))
	for k := range resp.Header {
		headers = append(headers, k)
	}
	sort.Strings(headers)
	for _, k := range headers {
		fmt.Printf("%s: %v\n", k, resp.Header[k])
	}
	fmt.Println()
	body, err := io.ReadAll(resp.Body) //TODO use tee() to be able to read the body again
	if err != nil {
		fmt.Println("error reading body!!!")
		return
	}
	fmt.Println(string(body))
}
