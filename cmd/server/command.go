package main

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/olahol/melody"
	log "github.com/sirupsen/logrus"
	errors "golang.org/x/xerrors"

	"github.com/carfloresf/financial-chat/config"
	"github.com/carfloresf/financial-chat/internal/controller"
	"github.com/carfloresf/financial-chat/internal/queue"
	"github.com/carfloresf/financial-chat/internal/service/user"
	"github.com/carfloresf/financial-chat/internal/storage"
)

// Execute starts the server, trying to keep all the logic related to starting the server in this function.
// nolint: cyclop
func Execute(configFile string) error {
	// Load configuration
	configuration, err := config.NewConfig(configFile)
	if err != nil {
		log.Errorf("error loading config: %s", err)

		return err
	}

	// start db connection
	storageDB, err := storage.NewStorage(&configuration.DB)
	if err != nil {
		log.Errorf("error creating storage: %s", err)

		return err
	}

	driver, err := sqlite3.WithInstance(storageDB.Conn, &sqlite3.Config{})
	if err != nil {
		log.Errorf("error creating driver: %s", err)

		return err
	}

	// start migration struct
	migration, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"main", driver)
	if err != nil {
		log.Errorf("migration error: %s", err)

		return err
	}

	// run migrations (in case something changed)
	if err = migration.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Errorf("migration up error: %s", err)

		return err
	}

	ph := user.NewPasswordGenerator(configuration.Auth.Pepper)

	userService := user.NewUser(ph, storageDB)

	// start queue client for publish, consume and close
	queueClient, err := queue.NewClient(configuration)
	if err != nil {
		log.Errorf("error creating queue: %s", err)
	}

	defer queueClient.Close()

	// start websocket server
	websocketServer := melody.New()
	queueClient.Websocket = websocketServer

	router, err := controller.NewRouter(configuration, userService, queueClient, websocketServer)
	if err != nil {
		log.Errorf("error creating router: %s", err)

		return err
	}

	// start response rabbitmq consumer
	err = queueClient.Consume(
		"chat_server",
		"response_queue",
		"response",
		queueClient.BroadcastAction)
	if err != nil {
		log.Errorf("error consuming queue: %s", err)

		return err
	}

	err = router.Run(configuration.HTTP.Addr + ":" + configuration.HTTP.Port)
	if err != nil {
		log.Errorf("failed to start server: %v", err)

		return err
	}

	return nil
}
