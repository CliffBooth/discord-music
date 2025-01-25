# DISOCRD - MUSIC

idea: write a discord bot that allows to play music (or any audio) from youtube by a link

[] disocrd bot docs: https://discord.com/developers/docs/quick-start/overview-of-apps

[] do we need configs?

---

## todo:
- [ ] create a simple prototype - hello world discord bot
    - [ ] research possibility of using discord library (better to do without it)
        https://github.com/bwmarrin/discordgo
        https://github.com/amatsagu/tempest

- [ ] take care of the rate limit! (research what it is and how not to get blocked)


- [ ] research library for downloading audio from youtube.
- [ ] research how you can actually stream audio (upload by chunks)

- [ ] run it in docker

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