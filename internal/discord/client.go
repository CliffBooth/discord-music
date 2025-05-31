package discord

import (
	"discord-music/internal/config" // TODO: if we want to treat this package as an external library we need to get rid of these imports
	"discord-music/internal/middleware"
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
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

	BOT_INTENTS = GUILDS | GUILD_MESSAGES /* | MESSAGE_CONTENT */ | DIRECT_MESSAGES | GUILD_VOICE_STATES
)

type Client struct {
	logger            log.Logger
	cfg               *config.Config
	httpClient        *http.Client
	commands          map[string]CommandFunc
	gateWayConnection *websocket.Conn
	wg                *sync.WaitGroup
	shutdownCh        chan struct{}

	state state
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
		commands:   map[string]CommandFunc{},
		wg:         &sync.WaitGroup{},
		shutdownCh: make(chan struct{}),
	}
}

type state struct {
	SequenceNumber  atomic.Int64
	SequenceStarted atomic.Bool // indicates whether sequenceNumber supposed to be nil
}

func (c *Client) Shutdown() {
	fmt.Println("inside shutdown")

	close(c.shutdownCh)

	if c.gateWayConnection != nil {
		c.gateWayConnection.Close()
	}

	c.wg.Wait()
	fmt.Println("inside shutdown after Wait()")
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

	c.gateWayConnection, _, err = websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		c.logger.Errorf("error = %v", err)
		return err
	}

	helloEvent := &HelloMessage{}
	err = c.gateWayConnection.ReadJSON(helloEvent)
	if err != nil {
		c.logger.Errorf("RunWebsocket() error: ", err)
		return err
	}
	c.logger.Debugf("helloEvent = %s", spew.Sprintf("%v", helloEvent))

	if helloEvent.S != nil {
		c.state.SequenceNumber.Store(*helloEvent.S)
		c.state.SequenceStarted.Store(true)
	}

	go c.heartbeatProc(c.gateWayConnection, helloEvent.D.HeartBeatInterval)

	//TODO: how to handle connection close -> resume ?
	// TODO: handle heartbeat requests - close current conenction and resume new one
	go c.websocketListen(c.gateWayConnection)

	message := IdentifyMessage{
		BaseMessage: BaseMessage{OP_IDENTIFY, nil},
		D: IdentifyData{
			Token:   c.cfg.TOKEN,
			Intents: BOT_INTENTS,
		},
	}
	err = c.gateWayConnection.WriteJSON(message)
	if err != nil {
		c.logger.Errorf("error sending identify message: %v", err)
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	c.Shutdown()

	return nil
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

// TODO: add jitter (as in docs)
// add checking whether reply has been received
func (c *Client) heartbeatProc(conn *websocket.Conn, interval int) {
	c.wg.Add(1)
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)

	for {
		select {
		case <-ticker.C:
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
				c.logger.Errorf("heartbeatProc() err: %v", err)
				return
			}
			c.logger.Debugf("heartbeat sent: %s", spew.Sprintf("%v", heartBeat))
		case <-c.shutdownCh:
			ticker.Stop()
			c.logger.Info("quit heartbeatProc")
			c.wg.Done()
			return
		}
	}
}

// TODO: need shutdown here
func (c *Client) websocketListen(conn *websocket.Conn) {
	c.wg.Add(1)

	errorHappened := false

	for {
		select {
		case <-c.shutdownCh:
			c.logger.Info("quit websocketListen")
			c.wg.Done()
			return

		default:
			if errorHappened {
				time.Sleep(time.Second)
				errorHappened = false
			}

			msgType, message, err := conn.ReadMessage()
			if err != nil {
				c.logger.Errorf("ReadMessage error: %v", err)
				errorHappened = true
				continue
			}

			if msgType == 1 {
				base := &BaseMessage{}
				err := json.Unmarshal(message, base)
				if err != nil {
					c.logger.Errorf("error unmarshalling base message: %v", err)
					continue
				}
				if base.S != nil {
					c.state.SequenceNumber.Store(*base.S)
				}
				c.logger.Debugf("received msg = %v", string(message))
				switch base.Op {
				case OP_DISPATCH:
					err := c.routeEvent(message, conn)
					if err != nil {
						c.logger.Errorf("websocketListen error: %v", err)
					}
				case OP_HEARTBEAT_ACK:
					heartbeatMsg := &HeartBeat{}
					err := json.Unmarshal(message, heartbeatMsg)
					if err != nil {
						c.logger.Errorf("error unmarshalling op=%d message: %v", base.Op, err)
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
				c.logger.Debugf("message type is not 1, but %d, msg=%s", msgType, string(message))
			}

		}
	}
}

func (c *Client) sendGateWayEvent(event any) error {
	if c.gateWayConnection == nil {
		return fmt.Errorf("no connection")
	}

	return c.gateWayConnection.WriteJSON(event)
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
	// case "VOICE_STATE_UPDATE":
	// 	return c.onVoiceStateUpdate(message)
	default:
		c.logger.Debugf("routeEvent: unknown T: %s", base.T)
		return nil
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

	commandText := event.D.Data.Name
	f, ok := c.commands[commandText]
	if !ok {
		c.logger.Info("unknown command: %s", commandText)
		return nil
	}

	err = f(&CommandContext{
		interactionID:    event.D.ID,
		interactionToken: event.D.Token,

		Event:  event,
		client: c,
	})
	if err != nil {
		c.logger.Errorf("command: %s, err: %e", commandText, err)
		return err
	}

	return nil
}

// this is how you add new commands to this discord client.
func (c *Client) OnCommand(text string, command CommandFunc) {
	c.commands[text] = command
}
