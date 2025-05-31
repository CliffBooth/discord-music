package discord

import "fmt"

type CommandContext struct {
	interactionID    string
	interactionToken string

	client *Client
	Event  *InteractionCreateEvent
}

type CommandFunc func(*CommandContext) error

var (
	ErrUserNotInVoiceChannel = fmt.Errorf("user not in voice channel")
)

// юзер нашей либы работает с CommandContext. Мы не даем ему использовать любой метод из CommandContext.client, поэтому деалем делегирование.
// мы не хотим например чтобы была возможность вызвать RunWebSocket()
func (c *CommandContext) SendText(content string) {
	response := InteractionResponse{
		Type: CHANNEL_MESSAGE_WITH_SOURCE,
		Data: InteractionRespData{
			Content: content,
		},
	}
	c.client.SendInteractionResponse(c.interactionID, c.interactionToken, response)
}

func (c *CommandContext) JoinVoiceChannel() error {
	state, err := c.client.GetUserVoiceState(c.Event.D.GuildID, c.Event.D.Member.User.ID)
	if err != nil {
		c.client.logger.Errorf("GetUserVoiceState() error: %v", err)
		return err
	}

	if state.ChannelID == "" {
		return ErrUserNotInVoiceChannel
	}

	voiceUpdateEvent := VoiceStateUpdateRequest{
		OP: 4,
		D: VoiceStateUpdateData{
			GuildID:   c.Event.D.GuildID,
			ChannelID: state.ChannelID,
		},
	}

	err = c.client.sendGateWayEvent(voiceUpdateEvent)
	if err != nil {
		c.client.logger.Errorf("JoinVoiceChannel() error: %v", err)
		return err
	}

	return nil
}

// func (c *CommandContext) PlayAudio() {

// }
