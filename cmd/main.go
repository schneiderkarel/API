package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"

	"app/cmd/config"
	"app/internal/helpers"
	"app/internal/httpserver"
	"app/internal/storage"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Init()

	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("cannot open postgres connection: error - %s", err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("cannot close postgres connection: error - %s", err.Error())
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatalf("cannot ping postgres connection: error - %s", err.Error())
	}

	httpServer := httpserver.New(
		cfg.HttpServerPort,
		storage.NewUserStorage(db),
	)

	httpServerErrCh := make(chan error, 1)
	go func() {
		httpServerErrCh <- httpServer.Run()
	}()

	systemSignalCh := make(chan os.Signal, 1)
	signal.Notify(systemSignalCh, os.Interrupt)
	var shutdownFn func()
	select {
	case <-systemSignalCh:
		shutdownFn = func() {
			httpServer.Stop()
		}
	case err := <-httpServerErrCh:
		shutdownFn = func() {
			log.Println(fmt.Errorf("http server unexpectedly stopped: %w", err))
		}
	}
	if ok := helpers.WithTimeout(shutdownFn, helpers.DefaultTimeout); !ok {
		log.Fatalln("graceful shutdown timed out")
	}

	log.Println("shut down")
}
