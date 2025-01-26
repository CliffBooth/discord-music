package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"discord-music/internal/config"
	"discord-music/internal/discord"
)

func main() {
	writers := []io.Writer{os.Stdout}

	file, err := os.OpenFile("bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		writers = append(writers, file)
	} else {
		fmt.Printf("error creating log file: %v\n", err)
	}

	cfg := config.Load()

	logger := log.New(io.MultiWriter(writers...), "", log.LstdFlags|log.Lshortfile)

	logger.Println("started")

	client := discord.New(logger, cfg)
	// client.ListCommands()
	// client.GetCurrentApplication()
	// client.InstallCommands()
	// discord.InstallDefualtCommands(client)
	client.RunWebsocket()
}
