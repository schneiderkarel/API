package httpserver

import (
	"app/internal/configuration"
	"app/internal/errors"
	"app/internal/storage"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	HeaderContentType     string = "Content-Type"
	HeaderContentTypeJson string = "application/json"
)

type HandlerInterface interface {
	HandleUsers(w http.ResponseWriter, r *http.Request)
	HandleUser(w http.ResponseWriter, r *http.Request)
	HandleCreateUser(w http.ResponseWriter, r *http.Request)
	HandleUpdateUser(w http.ResponseWriter, r *http.Request)
	HandleDeleteUser(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	ust storage.UserStorage
}

func NewHandler(db storage.UserStorage) HandlerInterface {
	return &handler{
		ust: db,
	}
}

func (h *handler) HandleUsers(w http.ResponseWriter, _ *http.Request) {
	users, err := h.ust.Users()
	if err != nil {
		errors.WriteInternalServerError(err, w)
		return
	}

	WriteJson(http.StatusOK, storage.UsersResponse{Users: users}, w)
}

func (h *handler) HandleUser(w http.ResponseWriter, r *http.Request) {
	var rb configuration.UserIdentifierRequest

	ok := ParseRequestBody(w, r, &rb)
	if !ok {
		return
	}

	if vErrs := rb.Validate(); len(vErrs) > 0 {
		errors.WriteUnprocessableEntitiesError(vErrs, w)
		return
	}

	usr, err := h.ust.User(rb.UserId)
	if err != nil {
		if err == storage.UserNotFoundErr {
			errors.WriteNotFoundError(err.Error(), w)
			return
		}

		errors.WriteInternalServerError(err, w)
		return
	}

	WriteJson(http.StatusOK, usr, w)
}

func (h *handler) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var rb configuration.UserRequest

	ok := ParseRequestBody(w, r, &rb)
	if !ok {
		return
	}

	if vErrs := rb.Validate(); len(vErrs) > 0 {
		errors.WriteUnprocessableEntitiesError(vErrs, w)
		return
	}

	err := h.ust.CreateUser(storage.User{UserId: rb.UserId, Name: rb.Name, Age: rb.Age})
	if err != nil {
		if err == storage.UserAlreadyExistsErr {
			errors.WriteConflictError(err.Error(), w)
			return
		}

		errors.WriteInternalServerError(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var rb configuration.UserRequest

	ok := ParseRequestBody(w, r, &rb)
	if !ok {
		return
	}

	if vErrs := rb.Validate(); len(vErrs) > 0 {
		errors.WriteUnprocessableEntitiesError(vErrs, w)
		return
	}

	err := h.ust.UpdateUser(storage.User{UserId: rb.UserId, Name: rb.Name, Age: rb.Age})
	if err != nil {
		if err == storage.UserNotFoundErr {
			errors.WriteNotFoundError(err.Error(), w)
			return
		}

		errors.WriteInternalServerError(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	var rb configuration.UserIdentifierRequest

	ok := ParseRequestBody(w, r, &rb)
	if !ok {
		return
	}

	if vErrs := rb.Validate(); len(vErrs) > 0 {
		errors.WriteUnprocessableEntitiesError(vErrs, w)
		return
	}

	err := h.ust.DeleteUser(rb.UserId)
	if err != nil {
		if err == storage.UserNotFoundErr {
			errors.WriteNotFoundError(err.Error(), w)
			return
		}

		errors.WriteInternalServerError(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ParseRequestBody(w http.ResponseWriter, r *http.Request, rb interface{}) bool {
	rawRb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errors.WriteInternalServerError(err, w)
		return false
	}

	if len(rawRb) == 0 {
		errors.WriteBadRequestError("empty request body", w)
		return false
	}

	if err = json.Unmarshal(rawRb, &rb); err != nil {
		errors.WriteBadRequestError("invalid request body", w)
		return false
	}

	return true
}

func WriteJson(statusCode int, body interface{}, w http.ResponseWriter) {
	content, err := json.Marshal(body)
	if err != nil {
		errors.WriteInternalServerError(err, w)
		return
	}

	w.Header().Set(HeaderContentType, HeaderContentTypeJson)
	w.WriteHeader(statusCode)

	if _, err := w.Write(content); err != nil {
		errors.WriteInternalServerError(err, w)
		log.Printf("unable to write http response - error: %s", err.Error())
	}
}
