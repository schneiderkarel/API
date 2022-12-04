package httpserver

import (
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter(h HandlerInterface) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/users", h.Users).Methods(http.MethodPost)
	router.HandleFunc("/user", h.User).Methods(http.MethodPost)
	router.HandleFunc("/create-user", h.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/update-user", h.UpdateUser).Methods(http.MethodPost)
	router.HandleFunc("/delete-user", h.DeleteUser).Methods(http.MethodPost)

	return router
}
