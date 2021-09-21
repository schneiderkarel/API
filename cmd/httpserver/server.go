package httpserver

import (
	"app/internal/storage"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const shutdownTimeout = 5 * time.Second

func Run(port int, userStorage storage.UserStorage) {
	defer func() {
		if r := recover(); r != nil {
			log.Fatalf("http server: fatal error - %e", r)
		}
	}()

	done := make(chan struct{})
	quit := make(chan os.Signal, 1)

	httpServer := &http.Server{
		Handler: NewRouter(
			NewHandler(userStorage),
		),
		Addr: fmt.Sprintf(":%d", port),
	}

	go gracefulShutdown(httpServer, quit, done)

	log.Printf("http server: started, port: %d", port)

	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Printf("http server: closed with error - %s", err.Error())
	}

	<-done
	log.Printf("http server: stopped, port: %d", port)
}

func gracefulShutdown(httpServer *http.Server, quit <-chan os.Signal, done chan<- struct{}) {
	<-quit
	log.Printf("http server: start graceful shutdown")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Printf("http server: shutdown error - %s", err.Error())
	}

	close(done)
}
