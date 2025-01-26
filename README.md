# DISOCRD - MUSIC

idea: write a discord bot that allows to play music (or any audio) from youtube by a link

[] disocrd bot docs: https://discord.com/developers/docs/quick-start/overview-of-apps

[x] do we need configs? - of course

---

## todo:
- [ ] create a simple prototype - hello world discord bot
    - [ ] research possibility of using discord library (better to do without it)
        https://github.com/bwmarrin/discordgo
        https://github.com/amatsagu/tempest

- [ ] take care of the rate limit! (research what it is and how not to get blocked)
    - [x] research


- [ ] research library for downloading audio from youtube.
- [ ] research how you can actually stream audio (upload by chunks)

- [ ] run it in docker

- [x] implement request middleware infrastructure
- [ ] finish implemening ratelimit middleware

- [ ] switch to a normal logger which supports log levels and log format (etc json)
    - [ ] also logger should do ratation log (for when i run the app in docker)
    - [ ] maybe impelement logger with log levels yourself

- [ ] try opening websocket. Can we read events from it? (instead of using public url?)

- [ ] cresate some sort of state machine for websocket connection

- [ ] Eventually we need to create our own discord-websocket library and use it from the outside (we need to maximally decouple websocket stuff from the business logic). Use it something like this:

```go
    func main() { 
        websocketClient := discord.NewWebsocketClient()
        websocket.OnCommand("test", func(c discord.MessageContext){}) // OnCommand adds a callback to when someone sends this command to the bot
        websocket.OnCommand("play", onPlay)
        websocket.Run() //this handles all the init stuff, hearbeat and reconnect stuff and so on, and also blocks.
    }

    func onPlay(c discord.MessageContext){...} // context has functions like c.Reply(string), c.SendSomething(), c.Speak(stream) - to play music
```

## planned features:
- [ ] add list of greetings, the bot will say a random one to greet you
- [ ] For extensiablity, allow http client in discord to add middlewares to it. Implement things like rate-limit checking in middleware. So it can look like this:

    ```go
    func NewDiscord() {
        httpClient := http.NewCLient()
        addMiddleware(httpClient, &RateLimitMiddleware{}) // when implemented a new feature, just add it in the Discord() constructor
    }

    func addMiddleware(c *http.Client, f func(next func() (http.Response, error)) error)

    func (r *RateLimitMiddleware) Apply(next func() (http.Response, error)) error {
        ok := check_ratelimit() // check cache
        if !ok {
            return ErrorRateLmit
        }
        resp, err := next()
        if err != nil {
            return err
        }
        update_ratelimit(resp) // update cache
    }
    ```
- [ ] add /мяу command which sends cat pictures 😺
- [ ] add rickroll somehow (maybe there will be a small chance that music will be rickroll instaead of requested music)
