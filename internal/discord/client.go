package discord

import (
	"discord-music/internal/config"
	httpClient "discord-music/internal/httpClient"
	httpclient "discord-music/internal/httpClient"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	BASE_URL = "https://discord.com/api/v10"
)

type Client struct {
	logger     *log.Logger
	cfg        *config.Config
	httpClient *httpclient.Client
}

func (c *Client) GetCurrentApplication() {
	endpoint := "/applications/@me"
	method := "GET"
	c.makeRequest(endpoint, method)
}

func (c *Client) makeRequest(
	endpoint string,
	method string,
) {
	url := BASE_URL + endpoint
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		c.logger.Printf("create request error: %v\n", err)
		return
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bot %s", c.cfg.TOKEN))
	request.Header.Add("User-Agent", "DiscordBot (https://github.com/CliffBooth/discord-music, 1.0.0)")
	request.Header.Add("Content-Type", "application/json; charset=UTF-8")
	// printRequest(request)
	fmt.Println()

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.logger.Printf("request error: %v\n", err)
		return
	}

	_ = response
	// printResponse(response)
}

func (c *Client) ListCommands() {
	endpoint := fmt.Sprintf("/applications/%s/commands", c.cfg.APP_ID)
	c.makeRequest(endpoint, "GET")
}

func (c *Client) InstallCommands(commands []Command) {
	endpoint := fmt.Sprintf("/applications/%s/commands", c.cfg.APP_ID)
	c.makeRequest(endpoint, "PUT")
}

func printRequest(req *http.Request) {
	fmt.Printf("%s %s\n", req.Method, req.URL)
	fmt.Println()
	for k, v := range req.Header {
		fmt.Printf("%s: %v\n", k, v)
	}
}

func printResponse(resp *http.Response) {
	fmt.Println(resp.StatusCode)
	fmt.Println()
	for k, v := range resp.Header {
		fmt.Printf("%s: %v\n", k, v)
	}
	fmt.Println()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading body!!!")
		return
	}
	fmt.Println(string(body))
}

func logMiddleware(next httpClient.Middleware) httpClient.Middleware {
	return func(req *http.Request) (*http.Response, error) {
		printRequest(req)
		resp, err := next(req)
		if err == nil {
			printResponse(resp)
		} else {
			fmt.Printf("error occured: %v\n", err)
		}
		return resp, err
	}
}

func New(logger *log.Logger, cfg *config.Config) *Client {
	c := http.Client{}

	httpClient := httpclient.New(c) // TODO: create own client

	httpClient.ApplyMiddlewares(
		RateLimitMiddleware,
		logMiddleware,
	)

	return &Client{
		logger:     logger, // TODO: maybe wrap this logger into another logger, so that this module has some tag like [discord-client] message...
		cfg:        cfg,
		httpClient: httpClient,
	}
}
