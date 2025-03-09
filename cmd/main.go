package main

import (
	"discord-music/internal/config"
	"discord-music/internal/discord"
	"discord-music/internal/log"

)

func main() {
	cfg := config.Load()

	logger, err := log.NewZapLogger(cfg)
	if err != nil {
		log.Default().Fatalf("error creating logger: %v\n", err)
	}

	logger.Info("started")

	client := discord.New(logger, cfg)
	// client.ListCommands()
	// client.GetCurrentApplication()
	// client.InstallCommands()
	// discord.InstallDefualtCommands(client)
	client.RunWebsocket()
}
