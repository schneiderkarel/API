package httpserver

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter(h HandlerInterface) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/users", h.HandleUsers).Methods(http.MethodPost)
	router.HandleFunc("/user", h.HandleUser).Methods(http.MethodPost)
	router.HandleFunc("/create-user", h.HandleCreateUser).Methods(http.MethodPost)
	router.HandleFunc("/update-user", h.HandleUpdateUser).Methods(http.MethodPost)
	router.HandleFunc("/delete-user", h.HandleDeleteUser).Methods(http.MethodPost)

	return router
}
