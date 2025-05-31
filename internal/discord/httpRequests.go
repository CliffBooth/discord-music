package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) makeRequest(
	endpoint string,
	method string,
	body io.Reader,
	headers http.Header,
) *http.Response {
	url := BASE_URL + endpoint
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		c.logger.Errorf("create request error: %v", err)
		return nil
	}

	for k, values := range headers {
		for _, v := range values {
			request.Header.Add(k, v)
		}
	}
	request.Header.Add("Authorization", fmt.Sprintf("Bot %s", c.cfg.TOKEN))
	request.Header.Add("User-Agent", "DiscordBot (https://github.com/CliffBooth/discord-music, 1.0.0)")

	response, err := c.httpClient.Do(request)
	if err != nil {
		c.logger.Errorf("request error: %v", err)
		return nil
	}

	return response
}

// makeJsonRequest simply converts input data to json and adds corresponding content-type header.
// This is a separate mehtod because we potentially can have somehting like makeUrlEncodedRequest and so on
func (c *Client) makeJsonRequest(endpoint, method string, d interface{}) error {
	data, err := json.Marshal(d)
	if err != nil {
		return err
	}
	body := bytes.NewReader(data)
	c.makeRequest(endpoint, method, body, http.Header(map[string][]string{"Content-Type": {"application/json; charset=UTF-8"}}))
	return nil
}

func (c *Client) GetCurrentApplication() {
	endpoint := "/applications/@me"
	method := "GET"
	c.makeRequest(endpoint, method, nil, nil)
}

func (c *Client) ListCommands() {
	endpoint := fmt.Sprintf("/applications/%s/commands", c.cfg.APP_ID)
	c.makeRequest(endpoint, http.MethodGet, nil, nil)
}

func (c *Client) InstallCommands(commands []Command) {
	endpoint := fmt.Sprintf("/applications/%s/commands", c.cfg.APP_ID)
	err := c.makeJsonRequest(endpoint, http.MethodPut, commands)
	if err != nil {
		c.logger.Errorf("InstallCommands error: %v", err)
	}
}

func (c *Client) GetGetaway() (string, error) {
	resp := &struct {
		Url string `json:"url"`
	}{}
	endpoint := "/gateway/bot"
	r := c.makeRequest(endpoint, http.MethodGet, nil, nil)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		c.logger.Errorf("GetGetaway() error reading response body: %v", err)
		return "", err
	}
	json.Unmarshal(body, resp)
	return resp.Url, nil
}

func (c *Client) SendInteractionResponse(interaction_id, interaction_token string, data InteractionResponse) {
	endpoint := fmt.Sprintf("/interactions/%s/%s/callback", interaction_id, interaction_token)
	c.makeJsonRequest(endpoint, http.MethodPost, data)
}

func (c *Client) GetUserVoiceState(guildID, userID string) (UserVoiceStateResponse, error) {
	endpoint := fmt.Sprintf("/guilds/%s/voice-states/%s", guildID, userID)
	response := c.makeRequest(endpoint, http.MethodGet, nil, nil)

	res := UserVoiceStateResponse{}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return res, err
	}

	json.Unmarshal(body, &res)
	return res, nil
}
