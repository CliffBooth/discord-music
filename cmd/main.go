package main

import (
	"discord-music/internal/config"
	"discord-music/internal/discord"
	"discord-music/internal/handlers"
	"discord-music/internal/log"
	"sync"
)

func main() {
	cfg := config.Load()

	logger, err := log.NewZapLogger(cfg)
	if err != nil {
		log.Default().Fatalf("error creating logger: %v\n", err)
	}

	logger.Info("started")

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func(wg *sync.WaitGroup) {
		client := discord.New(logger, cfg)
		// client.ListCommands()
		handlers.InstallDefualtCommands(client)
		err := client.RunWebsocket()
		if err != nil {
			logger.Errorf("runWebsocket error: %v", err)
		}
		wg.Done()
	}(wg)

	wg.Wait()

	logger.Info("quit main")
}
