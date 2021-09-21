package httpserver

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	expResponseBodyUsers      = "users OK"
	expResponseBodyUser       = "user OK"
	expResponseBodyCreateUser = "create-user OK"
	expResponseBodyUpdateUser = "update-user OK"
	expResponseBodyDeleteUser = "delete-user OK"
)

func TestNewRouter(t *testing.T) {
	bh := &baseHandlerMock{}
	router := NewRouter(bh)

	type args struct {
		method string
		url    string
	}
	type exp struct {
		respBody string
	}
	okTcs := []struct {
		name string
		args args
		exp  exp
	}{
		{
			name: "users",
			args: args{
				method: http.MethodPost,
				url:    "/users",
			},
			exp: exp{
				respBody: expResponseBodyUsers,
			},
		},
		{
			name: "user",
			args: args{
				method: http.MethodPost,
				url:    "/user",
			},
			exp: exp{
				respBody: expResponseBodyUser,
			},
		},
		{
			name: "create-user",
			args: args{
				method: http.MethodPost,
				url:    "/create-user",
			},
			exp: exp{
				respBody: expResponseBodyCreateUser,
			},
		},
		{
			name: "update-user",
			args: args{
				method: http.MethodPost,
				url:    "/update-user",
			},
			exp: exp{
				respBody: expResponseBodyUpdateUser,
			},
		},
		{
			name: "delete-user",
			args: args{
				method: http.MethodPost,
				url:    "/delete-user",
			},
			exp: exp{
				respBody: expResponseBodyDeleteUser,
			},
		},
	}

	for _, tc := range okTcs {
		t.Run("ok - "+tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.args.method, tc.args.url, bytes.NewReader([]byte("")))
			assert.NoError(t, err)

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, http.StatusOK, rec.Code, "unexpected response code")
			assert.Equal(t, tc.exp.respBody, rec.Body.String(), "unexpected response body")
		})
	}

	t.Run("not found", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/not-found", bytes.NewReader([]byte("")))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code, "unexpected response code")
	})

	t.Run("method not allowed", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPut, "/users", bytes.NewReader([]byte("")))
		assert.NoError(t, err)

		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code, "unexpected response code")
	})
}

type baseHandlerMock struct{}

func (bh *baseHandlerMock) HandleUsers(w http.ResponseWriter, r *http.Request) {
	bh.write(w, expResponseBodyUsers)
}

func (bh *baseHandlerMock) HandleUser(w http.ResponseWriter, r *http.Request) {
	bh.write(w, expResponseBodyUser)
}

func (bh *baseHandlerMock) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	bh.write(w, expResponseBodyCreateUser)
}

func (bh *baseHandlerMock) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	bh.write(w, expResponseBodyUpdateUser)
}

func (bh *baseHandlerMock) HandleDeleteUser(w http.ResponseWriter, r *http.Request) {
	bh.write(w, expResponseBodyDeleteUser)
}

func (bh *baseHandlerMock) write(w http.ResponseWriter, responseBody string) {
	_, err := w.Write([]byte(responseBody))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
