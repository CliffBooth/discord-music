package handlers

import (
	"discord-music/internal/discord"
	"fmt"
)

var _ discord.CommandFunc = TestCommand

func TestCommand(ctx *discord.CommandContext) error {

	username := ctx.Event.D.Member.User.Username
	responseMsg := fmt.Sprintf("hello, %s 😸", username)

	ctx.SendText(responseMsg)

	return nil
}
