package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/carfloresf/financial-chat/config"
	"github.com/carfloresf/financial-chat/internal/queue"
	"github.com/carfloresf/financial-chat/internal/service/bot"
	"github.com/carfloresf/financial-chat/internal/stooq"
)

func Execute(configFile string) error {
	configuration, err := config.NewConfig(configFile)
	if err != nil {
		log.Errorf("error loading config: %s", err)

		return err
	}

	stooqClient := stooq.NewClient()

	queueClient, err := queue.NewClient(configuration)
	if err != nil {
		log.Errorf("error creating queue: %s", err)

		return err
	}

	defer queueClient.Close()

	commandExecutor := bot.NewExecutor(stooqClient, queueClient)

	go func() {
		err = queueClient.Consume(
			"bot",
			"request_queue",
			"request",
			commandExecutor.Execute)
		if err != nil {
			log.Errorf("error starting consumer: %s", err)

			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	return nil
}
