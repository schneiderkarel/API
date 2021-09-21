package main

import (
	"app/cmd/config"
	"app/cmd/httpserver"
	"app/internal/storage"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"sync"
)

func main() {
	cfg := config.Init()

	db, err := sql.Open("postgres", cfg.PostgresDSN)
	if err != nil {
		log.Printf("cannot open postgres connection: error - %s", err.Error())
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("cannot close postgres connection: error - %s", err.Error())
		}
	}()

	if err := db.Ping(); err != nil {
		log.Printf("cannot ping postgres connection: error - %s", err.Error())
	}

	userStorage := storage.NewUserStorage(db)

	wg := &sync.WaitGroup{}

	addToWaitGroup(wg, func() {
		httpserver.Run(cfg.HttpServerPort, userStorage)
	})

	wg.Wait()
}

func addToWaitGroup(wg *sync.WaitGroup, f func()) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		f()
	}()
}
