package handlers

import (
	"discord-music/internal/discord"
	"errors"
)

var _ discord.CommandFunc = PlayCommand

// здесь есть как напрмер замутить другого юзера https://discord.com/developers/docs/resources/voice#voice-state-object

func PlayCommand(ctx *discord.CommandContext) error {
	fileaname := "01-Me And Your Mama.mp3"

	err := ctx.JoinVoiceChannel()
	switch {
	case errors.Is(discord.ErrUserNotInVoiceChannel, err):
		ctx.SendText("bro, you are not even in a voice channel, the fuck you want from me...")
	case err == nil:
	default:
		ctx.SendText("something went wrong, try again...")
		return err
	}
	ctx.SendText("joining the voice channel...")

	ctx.SendText("playing: " + fileaname)
	return nil
}
