package httpserver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"app/internal/configuration"
	"app/internal/response"
	"app/internal/storage"
)

type HandlerInterface interface {
	Users(w http.ResponseWriter, r *http.Request)
	User(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
}

type handler struct {
	ust storage.UserStorage
}

func newHandler(ust storage.UserStorage) HandlerInterface {
	return &handler{
		ust: ust,
	}
}

func (h *handler) Users(w http.ResponseWriter, _ *http.Request) {
	users, err := h.ust.Users()
	if err != nil {
		response.WriteInternalServerError(err, w)
		return
	}

	response.WriteJson(http.StatusOK, storage.UsersResponse{Users: users}, w)
}

func (h *handler) User(w http.ResponseWriter, r *http.Request) {
	var rb configuration.UserIdentifierRequest

	ok := parseRequestBody(w, r, &rb)
	if !ok {
		return
	}

	if vErrs := rb.Validate(); len(vErrs) > 0 {
		response.WriteUnprocessableEntitiesError(vErrs, w)
		return
	}

	usr, err := h.ust.User(rb.UserId)
	if err != nil {
		if err == storage.UserNotFoundErr {
			response.WriteNotFoundError(err.Error(), w)
			return
		}

		response.WriteInternalServerError(err, w)
		return
	}

	response.WriteJson(http.StatusOK, usr, w)
}

func (h *handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var rb configuration.UserRequest

	ok := parseRequestBody(w, r, &rb)
	if !ok {
		return
	}

	if vErrs := rb.Validate(); len(vErrs) > 0 {
		response.WriteUnprocessableEntitiesError(vErrs, w)
		return
	}

	err := h.ust.CreateUser(storage.User{UserId: rb.UserId, Name: rb.Name, Age: rb.Age})
	if err != nil {
		if err == storage.UserAlreadyExistsErr {
			response.WriteConflictError(err.Error(), w)
			return
		}

		response.WriteInternalServerError(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var rb configuration.UserRequest

	ok := parseRequestBody(w, r, &rb)
	if !ok {
		return
	}

	if vErrs := rb.Validate(); len(vErrs) > 0 {
		response.WriteUnprocessableEntitiesError(vErrs, w)
		return
	}

	err := h.ust.UpdateUser(storage.User{UserId: rb.UserId, Name: rb.Name, Age: rb.Age})
	if err != nil {
		if err == storage.UserNotFoundErr {
			response.WriteNotFoundError(err.Error(), w)
			return
		}

		response.WriteInternalServerError(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var rb configuration.UserIdentifierRequest

	ok := parseRequestBody(w, r, &rb)
	if !ok {
		return
	}

	if vErrs := rb.Validate(); len(vErrs) > 0 {
		response.WriteUnprocessableEntitiesError(vErrs, w)
		return
	}

	err := h.ust.DeleteUser(rb.UserId)
	if err != nil {
		if err == storage.UserNotFoundErr {
			response.WriteNotFoundError(err.Error(), w)
			return
		}

		response.WriteInternalServerError(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func parseRequestBody(w http.ResponseWriter, r *http.Request, rb interface{}) bool {
	rawRb, err := ioutil.ReadAll(r.Body)
	if err != nil {
		response.WriteInternalServerError(err, w)
		return false
	}

	if len(rawRb) == 0 {
		response.WriteBadRequestError("empty request body", w)
		return false
	}

	if err = json.Unmarshal(rawRb, &rb); err != nil {
		response.WriteBadRequestError("invalid request body", w)
		return false
	}

	return true
}
