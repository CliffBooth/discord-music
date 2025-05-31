package handlers

import "discord-music/internal/discord"

var (
	defaultCommands = []discord.Command{
		{
			Name:              "test",
			Description:       "hello biatch!!!",
			Type:              1,
			Integration_types: []int{0, 1},
			Contexts:          []int{0, 1, 2},
			F:                 TestCommand,
		},
		{
			Name:              "play1",
			Description:       "play some music",
			Type:              1,
			Integration_types: []int{0, 1},
			Contexts:          []int{0, 1, 2},
			F:                 PlayCommand,
		},
	}
)

func InstallDefualtCommands(client *discord.Client) {
	client.InstallCommands(defaultCommands)

	for _, c := range defaultCommands {
		client.OnCommand(c.Name, c.F)
	}
}
