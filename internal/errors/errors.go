package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

const (
	headerContentType     = "Content-Type"
	headerContentTypeJson = "application/json"
)

type badRequestError struct {
	Error string `json:"error"`
}

func WriteBadRequestError(err string, w http.ResponseWriter) {
	WriteJson(http.StatusBadRequest, badRequestError{Error: err}, w)
}

type notFoundError struct {
	Error string `json:"error"`
}

func WriteNotFoundError(err string, w http.ResponseWriter) {
	WriteJson(http.StatusNotFound, notFoundError{Error: err}, w)
}

type conflictError struct {
	Error string `json:"error"`
}

func WriteConflictError(err string, w http.ResponseWriter) {
	WriteJson(http.StatusConflict, conflictError{Error: err}, w)
}

type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

func WriteUnprocessableEntitiesError(vErrs []ValidationError, w http.ResponseWriter) {
	WriteJson(http.StatusUnprocessableEntity, ValidationErrors{Errors: vErrs}, w)
}

func WriteInternalServerError(err error, w http.ResponseWriter) {
	log.Printf("internal server error - %s", err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func WriteJson(statusCode int, body interface{}, w http.ResponseWriter) {
	content, err := json.Marshal(body)
	if err != nil {
		WriteInternalServerError(err, w)
		return
	}

	w.Header().Set(headerContentType, headerContentTypeJson)
	w.WriteHeader(statusCode)

	if _, err := w.Write(content); err != nil {
		WriteInternalServerError(err, w)
		log.Printf("unable to write http response - error: %s", err.Error())
	}
}
