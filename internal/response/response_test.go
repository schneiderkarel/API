package response

import (
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_WriteBadRequestError(t *testing.T) {
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteBadRequestError("bad-request", w)
	})

	req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(""))
	require.NoError(t, err)

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestWriteNotFoundError(t *testing.T) {
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteNotFoundError("not-found", w)
	})

	req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(""))
	require.NoError(t, err)

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusNotFound, res.Code)
}

func Test_WriteConflictError(t *testing.T) {
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteConflictError("conflict", w)
	})

	req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(""))
	require.NoError(t, err)

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusConflict, res.Code)
}

func Test_WriteUnprocessableEntitiesError(t *testing.T) {
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteUnprocessableEntitiesError([]ValidationError{}, w)
	})

	req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(""))
	require.NoError(t, err)

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnprocessableEntity, res.Code)
}

func Test_WriteInternalServerError(t *testing.T) {
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		WriteInternalServerError(errors.New("an error"), w)
	})

	req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(""))
	require.NoError(t, err)

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	assert.Equal(t, "", res.Body.String())
}

func Test_WriteJson(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		expCode := 499

		res := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WriteJson(
				expCode,
				struct {
					Status string `json:"status"`
				}{
					Status: "error",
				},
				w,
			)
		})

		req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(""))
		require.NoError(t, err)

		handler.ServeHTTP(res, req)

		assert.Equal(t, expCode, res.Code, "unexpected code")
		assert.Equal(t, headerContentTypeJson, res.Header().Get(headerContentType), "unexpected content type")
		assert.Equal(t, `{"status":"error"}`, res.Body.String(), "unexpected body")
	})

	t.Run("marshal error", func(t *testing.T) {
		res := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WriteJson(
				499,
				struct {
					Status float64 `json:"status"`
				}{Status: math.Inf(1)},
				w,
			)
		})

		req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(""))
		require.NoError(t, err)

		handler.ServeHTTP(res, req)

		assert.Equal(t, http.StatusInternalServerError, res.Code, "unexpected code")
		assert.Equal(t, "", res.Header().Get(headerContentType), "unexpected content type")
		assert.Equal(t, "", res.Body.String(), "unexpected body")
	})
}
