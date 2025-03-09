package discord

import (
	"bytes"
	"discord-music/internal/config"
	"discord-music/internal/middleware"
	"encoding/json"
	"io"
	"sync/atomic"
	"time"

	"fmt"
	"net/http"

	"discord-music/internal/log"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
)

const (
	BASE_URL = "https://discord.com/api/v10"
)

const (
	OP_DISPATCH              = 0
	OP_HEARTBEAT             = 1
	OP_IDENTIFY              = 2
	OP_PRESENCE_UPDATE       = 3
	OP_VOICE_STATE_UPDATE    = 4
	OP_RESUME                = 6
	OP_RECONNECT             = 7
	OP_REQUEST_GUILD_MEMBERS = 8
	OP_INVALID_SESSION       = 9
	OP_HELLO                 = 10
	OP_HEARTBEAT_ACK         = 11

	OP_REQUEST_SOUNDBOARD_SOUNDS = 31
)

// intents
const (
	GUILDS = 1 << iota
	GUILD_MEMBERS
	GUILD_MODERATION
	GUILD_EXPRESSIONS
	GUILD_INTEGRATIONS
	GUILD_WEBHOOKS
	GUILD_INVITES
	GUILD_VOICE_STATES
	GUILD_PRESENCES
	GUILD_MESSAGES
	GUILD_MESSAGE_REACTIONS
	GUILD_MESSAGE_TYPING
	DIRECT_MESSAGES
	DIRECT_MESSAGE_REACTIONS
	DIRECT_MESSAGE_TYPING
	MESSAGE_CONTENT
	GUILD_SCHEDULED_EVENTS

	AUTO_MODERATION_CONFIGURATION = 1 << 20
	AUTO_MODERATION_EXECUTION     = 1 << 21
	GUILD_MESSAGE_POLLS           = 1 << 24
	DIRECT_MESSAGE_POLLS          = 1 << 25

	BOT_INTENTS = GUILDS | GUILD_MESSAGES /* | MESSAGE_CONTENT */ | DIRECT_MESSAGES
)

type Client struct {
	logger     log.Logger
	cfg        *config.Config
	httpClient *http.Client

	state state
}

type state struct {
	SequenceNumber  atomic.Int64
	SequenceStarted atomic.Bool // indicates whether sequenceNumber supposed to be nil
}

func (c *Client) GetCurrentApplication() {
	endpoint := "/applications/@me"
	method := "GET"
	c.makeRequest(endpoint, method, nil, nil)
}

func (c *Client) makeRequest(
	endpoint string,
	method string,
	body io.Reader,
	headers http.Header,
) *http.Response {
	url := BASE_URL + endpoint
	request, err := http.NewRequest(method, url, body)
	if err != nil {
		c.logger.Errorf("create request error: %v\n", err)
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
		c.logger.Errorf("request error: %v\n", err)
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

func (c *Client) ListCommands() {
	endpoint := fmt.Sprintf("/applications/%s/commands", c.cfg.APP_ID)
	c.makeRequest(endpoint, http.MethodGet, nil, nil)
}

func (c *Client) InstallCommands(commands []Command) {
	endpoint := fmt.Sprintf("/applications/%s/commands", c.cfg.APP_ID)
	err := c.makeJsonRequest(endpoint, http.MethodPut, commands)
	if err != nil {
		c.logger.Errorf("InstallCommands error: %v\n", err)
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

type BaseMessage struct {
	Op int    `json:"op"`
	S  *int64 `json:"s"`
}

// TODO: later try to get rid of all these separate structs and just have all the possible fields in one struct with omitempty
type HelloMessage struct {
	BaseMessage
	D struct {
		HeartBeatInterval int `json:"heartbeat_interval"`
	} `json:"d"`
}

type HeartBeat struct {
	BaseMessage
	D *int64 `json:"d"`
}

type IdentifyMessage struct {
	BaseMessage
	D IdentifyData `json:"d"`
}

type IdentifyData struct {
	Token         string             `json:"token"`
	Intents       int                `json:"intents"`
	Properties    IdentifyProperties `json:"properties,omitempty"`
	Compress      *bool              `json:"compress,omitempty"`
	LargeTreshold *int               `json:"large_treshold,omitempty"`
	Shard         []int              `json:"shard,omitempty"`
	Presence      IdentifyPresence   `json:"presence,omitempty"`
}

type IdentifyProperties struct {
	Os      string `json:"os,omitempty"`
	Device  string `json:"device,omitempty"`
	Browser string `json:"browser,omitempty"`
}

type IdentifyPresence struct {
	//TODO
}

func (c *Client) RunWebsocket() error {
	// var upgrader = websocket.Upgrader{
	// 	ReadBufferSize:  1024,
	// 	WriteBufferSize: 1024,
	// }

	//TODO: cache value "url" from here
	url, err := c.GetGetaway()
	if err != nil {
		return err
	}
	c.logger.Debugf("websocket url = %s", url) // TODO: need to append ?v=10&encoding=json at the end.

	conn, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		c.logger.Errorf("error = %v\n", err)
		return err
	}
	defer conn.Close()

	helloEvent := &HelloMessage{}
	err = conn.ReadJSON(helloEvent)
	if err != nil {
		c.logger.Errorf("RunWebsocket() error: ", err)
		return err
	}
	c.logger.Debugf("helloEvent = %s", spew.Sprintf("%v", helloEvent))

	if helloEvent.S != nil {
		c.state.SequenceNumber.Store(*helloEvent.S)
		c.state.SequenceStarted.Store(true)
	}

	go c.heartbeatProc(conn, helloEvent.D.HeartBeatInterval)

	//TODO: how to handle connection close -> resume ?
	// TODO: handle heartbeat requests - close current conenction and resume new one
	go c.websocketListen(conn)

	message := IdentifyMessage{
		BaseMessage: BaseMessage{OP_IDENTIFY, nil},
		D: IdentifyData{
			Token:   c.cfg.TOKEN,
			Intents: BOT_INTENTS,
		},
	}
	err = conn.WriteJSON(message)
	if err != nil {
		c.logger.Errorf("error sending identify message: %v", err)
		return err
	}

	hearBeatStartedCh := make(chan struct{})
	<-hearBeatStartedCh //delete this later
	return nil
}

// TODO: add jitter (as in docs)
// add checking whether reply has been received
func (c *Client) heartbeatProc(conn *websocket.Conn, interval int) {
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	for range ticker.C {
		loaded := c.state.SequenceNumber.Load()
		var d = &loaded
		if c.state.SequenceStarted.Load() {
			d = nil
		}
		heartBeat := &HeartBeat{
			BaseMessage: BaseMessage{OP_HEARTBEAT, nil},
			D:           d,
		}
		err := conn.WriteJSON(heartBeat) //TODO: if don't receive ack from the server, must close this connection and Resume
		if err != nil {
			c.logger.Errorf("heartbeatProc() err: ", err)
			return
		}
		c.logger.Debugf("heartbeat sent: %s", spew.Sprintf("%v", heartBeat))
	}
}

func (c *Client) websocketListen(conn *websocket.Conn) {
	for {
		msgType, message, err := conn.ReadMessage()
		if err != nil {
			c.logger.Errorf("ReadMessage error: ", err)
			return
		}
		if msgType == 1 {
			base := &BaseMessage{}
			err := json.Unmarshal(message, base)
			if err != nil {
				c.logger.Errorf("error unmarshalling base message: ", err)
				continue
			}
			if base.S != nil {
				c.state.SequenceNumber.Store(*base.S)
			}
			c.logger.Debugf("received msg = %v\n", string(message))
			switch base.Op {
			case OP_DISPATCH:
				err := c.routeEvent(message, conn)
				if err != nil {
					c.logger.Errorf("websocketListen error: ", err)
				}
			case OP_HEARTBEAT_ACK:
				heartbeatMsg := &HeartBeat{}
				err := json.Unmarshal(message, heartbeatMsg)
				if err != nil {
					c.logger.Errorf("error unmarshalling op=%d message: %v\n", base.Op, err)
					continue
				}
				// if !heartBeatStarted {
				// 	heartBeatStarted = true
				// 	// value = heartbeatMsg.S
				// 	// hearBeatStartedCh <- struct{}{} // send command to the main process to proceed
				// }
			default:
				c.logger.Debugf("OPCODE RECEIVED: ", base.Op)
			}
		} else {
			c.logger.Debugf("message type is not 1, but %d, msg=%s\n", msgType, string(message))
		}
	}
}

func (c *Client) routeEvent(message []byte, conn *websocket.Conn) error {
	base := &BaseEvent{}
	err := json.Unmarshal(message, base)
	if err != nil {
		return fmt.Errorf("routeEvent: %v", err)
	}
	switch base.T {
	case "READY":
		return c.onReadyEvent(message, conn)
	case "GUILD_CREATE":
		return c.onGuildCreateEvent(message, conn)
	case "INTERACTION_CREATE":
		return c.onInteraction(message)
	default:
		return fmt.Errorf("routeEvent: unknown T: %s", base.T)
	}
}

func (c *Client) onReadyEvent(message []byte, conn *websocket.Conn) error {
	event := &ReadyEvent{}
	err := json.Unmarshal(message, event)
	if err != nil {
		c.logger.Infof("onReadyEvent() error: %v", err)
		return err
	}

	//TODO: save reconnect url

	return nil
}

func (c *Client) onGuildCreateEvent(message []byte, conn *websocket.Conn) error {
	event := &GuildCreateEvent{}
	err := json.Unmarshal(message, event)
	if err != nil {
		c.logger.Infof("onGuildCreateEvent() error: %v", err)
		return err
	}

	return nil
}

// https://discord.com/developers/docs/interactions/receiving-and-responding#interaction-callback
// If you are receiving Interactions over the gateway, you have to respond via HTTP. Responses to Interactions are not sent as commands over the gateway.
func (c *Client) onInteraction(message []byte) error {
	event := &InteractionCreateEvent{}
	err := json.Unmarshal(message, event)
	if err != nil {
		c.logger.Errorf("onInteraction() error: %v", err)
		return err
	}

	username := event.D.Member.User.Username
	responseMsg := fmt.Sprintf("hello, %s 😸", username)
	response := InteractionResponse{
		Type: CHANNEL_MESSAGE_WITH_SOURCE,
		Data: InteractionRespData{
			Content: responseMsg,
		},
	}

	c.SendInteractionResponse(event.D.ID, event.D.Token, response)

	return nil
}

func New(logger log.Logger, cfg *config.Config) *Client {
	httpClient := &http.Client{
		Transport: http.DefaultTransport,
	} //TODO: configure

	middleware.ApplyMiddlewares(
		httpClient,
		middleware.RateLimitMiddleware,
		middleware.GetLogMiddleware(logger),
	)

	return &Client{
		logger:     logger, // TODO: maybe wrap this logger into another logger, so that this module has some tag like [discord-client] message...
		cfg:        cfg,
		httpClient: httpClient,
		state:      state{},
	}
}
