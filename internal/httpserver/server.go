package httpserver

import (
	"fmt"
	"log"
	"net/http"

	"app/internal/storage"
)

type server struct {
	srv *http.Server
}

func New(
	port int,
	userStorage storage.UserStorage,
) *server {
	return &server{
		srv: &http.Server{
			Addr: fmt.Sprintf(":%d", port),
			Handler: newRouter(
				newHandler(
					userStorage,
				),
			),
		},
	}
}

func (s *server) Run() error {
	log.Println(fmt.Sprintf("http server listening on %s", s.srv.Addr))
	log.Println("http server successfully started")

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Println(fmt.Errorf("http server finished with an error - %s", err))
	} else {
		log.Println("http server run method finished")
	}

	return nil
}

func (s *server) Stop() {
	log.Println("http server stopping")

	if err := s.srv.Close(); err != nil {
		log.Println(fmt.Errorf("http server closing failed with error - %s", err))
	} else {
		log.Println("http server successfully stopped")
	}
}
