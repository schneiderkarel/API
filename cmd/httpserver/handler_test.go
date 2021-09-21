package httpserver

import (
	"app/internal/configuration"
	"app/internal/storage"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	user1 = storage.User{
		UserId: "7df661d5-47e3-4533-baa6-5f952d18bffe",
		Name:   "John Doe",
		Age:    42,
	}

	user2 = storage.User{
		UserId: "63df08d2-fa53-4575-a681-99058f8daba5",
		Name:   "Josh Brave",
		Age:    20,
	}
)

func TestNewHandler(t *testing.T) {
	expHandler := &handler{
		ust: userStorageMock{},
	}

	assert.Equal(t, expHandler, NewHandler(userStorageMock{}))
}

func TestHandler_HandleUsers(t *testing.T) {
	type args struct {
		ust storage.UserStorage
	}
	type exp struct {
		respCode int
		resBody  string
	}
	tcs := []struct {
		name string
		args args
		exp  exp
	}{
		{
			name: "ok",
			args: args{
				ust: userStorageMock{
					users: func() ([]storage.User, error) {
						return []storage.User{user1, user2}, nil
					},
				},
			},
			exp: exp{
				respCode: http.StatusOK,
				resBody:  `{"users":[{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe","name":"John Doe","age":42},{"user_id":"63df08d2-fa53-4575-a681-99058f8daba5","name":"Josh Brave","age":20}]}`,
			},
		},
		{
			name: "database error",
			args: args{
				ust: userStorageMock{
					users: func() ([]storage.User, error) {
						return nil, errors.New("database error")
					},
				},
			},
			exp: exp{
				respCode: http.StatusInternalServerError,
				resBody:  ``,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("", "", nil)
			require.NoError(t, err)

			rec := httptest.NewRecorder()

			h := NewHandler(tc.args.ust)

			h.HandleUsers(rec, req)

			assert.Equal(t, tc.exp.respCode, rec.Code, "unexpected response code")
			assert.Equal(t, tc.exp.resBody, rec.Body.String(), "unexpected response body")
		})
	}
}

func TestHandler_HandleUser(t *testing.T) {
	type args struct {
		reqBody string
		ust     storage.UserStorage
	}
	type exp struct {
		respCode int
		resBody  string
	}
	tcs := []struct {
		name string
		args args
		exp  exp
	}{
		{
			name: "ok",
			args: args{
				reqBody: `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe"}`,
				ust: userStorageMock{
					user: func() (storage.User, error) {
						return user1, nil
					},
				},
			},
			exp: exp{
				respCode: http.StatusOK,
				resBody:  `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe","name":"John Doe","age":42}`,
			},
		},
		{
			name: "invalid request body",
			args: args{
				reqBody: `{`,
				ust:     userStorageMock{},
			},
			exp: exp{
				respCode: http.StatusBadRequest,
				resBody:  `{"error":"invalid request body"}`,
			},
		},
		{
			name: "invalid values in request body",
			args: args{
				reqBody: `{"user_id":"u-1"}`,
				ust:     userStorageMock{},
			},
			exp: exp{
				respCode: http.StatusUnprocessableEntity,
				resBody:  `{"errors":[{"path":"user_id","message":"invalid UUID length: 3"}]}`,
			},
		},
		{
			name: "user not found error",
			args: args{
				reqBody: `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe"}`,
				ust: userStorageMock{
					user: func() (storage.User, error) {
						return storage.User{}, storage.UserNotFoundErr
					},
				},
			},
			exp: exp{
				respCode: http.StatusNotFound,
				resBody:  `{"error":"user not found"}`,
			},
		},
		{
			name: "database error",
			args: args{
				reqBody: `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe"}`,
				ust: userStorageMock{
					user: func() (storage.User, error) {
						return storage.User{}, errors.New("database error")
					},
				},
			},
			exp: exp{
				respCode: http.StatusInternalServerError,
				resBody:  ``,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("", "", strings.NewReader(tc.args.reqBody))
			require.NoError(t, err)

			rec := httptest.NewRecorder()

			h := NewHandler(tc.args.ust)

			h.HandleUser(rec, req)

			assert.Equal(t, tc.exp.respCode, rec.Code, "unexpected response code")
			assert.Equal(t, tc.exp.resBody, rec.Body.String(), "unexpected response body")
		})
	}
}

func TestHandler_HandleCreateUser(t *testing.T) {
	type args struct {
		reqBody string
		ust     storage.UserStorage
	}
	type exp struct {
		respCode int
		resBody  string
	}
	tcs := []struct {
		name string
		args args
		exp  exp
	}{
		{
			name: "ok",
			args: args{
				reqBody: `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe","name":"John Doe","age":42}`,
				ust: userStorageMock{
					createUser: func() error {
						return nil
					},
				},
			},
			exp: exp{
				respCode: http.StatusNoContent,
				resBody:  ``,
			},
		},
		{
			name: "invalid request body",
			args: args{
				reqBody: `{`,
				ust:     userStorageMock{},
			},
			exp: exp{
				respCode: http.StatusBadRequest,
				resBody:  `{"error":"invalid request body"}`,
			},
		},
		{
			name: "invalid values in request body",
			args: args{
				reqBody: `{"user_id":"u-1","name":"John Doe","age":42}`,
				ust:     userStorageMock{},
			},
			exp: exp{
				respCode: http.StatusUnprocessableEntity,
				resBody:  `{"errors":[{"path":"user_id","message":"invalid UUID length: 3"}]}`,
			},
		},
		{
			name: "user already exists error",
			args: args{
				reqBody: `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe","name":"John Doe","age":42}`,
				ust: userStorageMock{
					createUser: func() error {
						return storage.UserAlreadyExistsErr
					},
				},
			},
			exp: exp{
				respCode: http.StatusConflict,
				resBody:  `{"error":"user already exists"}`,
			},
		},
		{
			name: "database error",
			args: args{
				reqBody: `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe","name":"John Doe","age":42}`,
				ust: userStorageMock{
					createUser: func() error {
						return errors.New("database error")
					},
				},
			},
			exp: exp{
				respCode: http.StatusInternalServerError,
				resBody:  ``,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("", "", strings.NewReader(tc.args.reqBody))
			require.NoError(t, err)

			rec := httptest.NewRecorder()

			h := NewHandler(tc.args.ust)

			h.HandleCreateUser(rec, req)

			assert.Equal(t, tc.exp.respCode, rec.Code, "unexpected response code")
			assert.Equal(t, tc.exp.resBody, rec.Body.String(), "unexpected response body")
		})
	}
}

func TestHandler_HandleUpdateUser(t *testing.T) {
	type args struct {
		reqBody string
		ust     storage.UserStorage
	}
	type exp struct {
		respCode int
		resBody  string
	}
	tcs := []struct {
		name string
		args args
		exp  exp
	}{
		{
			name: "ok",
			args: args{
				reqBody: `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe","name":"John Doe","age":42}`,
				ust: userStorageMock{
					updateUser: func() error {
						return nil
					},
				},
			},
			exp: exp{
				respCode: http.StatusNoContent,
				resBody:  ``,
			},
		},
		{
			name: "invalid request body",
			args: args{
				reqBody: `{`,
				ust:     userStorageMock{},
			},
			exp: exp{
				respCode: http.StatusBadRequest,
				resBody:  `{"error":"invalid request body"}`,
			},
		},
		{
			name: "invalid values in request body",
			args: args{
				reqBody: `{"user_id":"u-1","name":"John Doe","age":42}`,
				ust:     userStorageMock{},
			},
			exp: exp{
				respCode: http.StatusUnprocessableEntity,
				resBody:  `{"errors":[{"path":"user_id","message":"invalid UUID length: 3"}]}`,
			},
		},
		{
			name: "user not found error",
			args: args{
				reqBody: `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe","name":"John Doe","age":42}`,
				ust: userStorageMock{
					updateUser: func() error {
						return storage.UserNotFoundErr
					},
				},
			},
			exp: exp{
				respCode: http.StatusNotFound,
				resBody:  `{"error":"user not found"}`,
			},
		},
		{
			name: "database error",
			args: args{
				reqBody: `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe","name":"John Doe","age":42}`,
				ust: userStorageMock{
					updateUser: func() error {
						return errors.New("database error")
					},
				},
			},
			exp: exp{
				respCode: http.StatusInternalServerError,
				resBody:  ``,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("", "", strings.NewReader(tc.args.reqBody))
			require.NoError(t, err)

			rec := httptest.NewRecorder()

			h := NewHandler(tc.args.ust)

			h.HandleUpdateUser(rec, req)

			assert.Equal(t, tc.exp.respCode, rec.Code, "unexpected response code")
			assert.Equal(t, tc.exp.resBody, rec.Body.String(), "unexpected response body")
		})
	}
}

func TestHandler_HandleDeleteUser(t *testing.T) {
	type args struct {
		reqBody string
		ust     storage.UserStorage
	}
	type exp struct {
		respCode int
		resBody  string
	}
	tcs := []struct {
		name string
		args args
		exp  exp
	}{
		{
			name: "ok",
			args: args{
				reqBody: `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe"}`,
				ust: userStorageMock{
					deleteUser: func() error {
						return nil
					},
				},
			},
			exp: exp{
				respCode: http.StatusNoContent,
				resBody:  ``,
			},
		},
		{
			name: "invalid request body",
			args: args{
				reqBody: `{`,
				ust:     userStorageMock{},
			},
			exp: exp{
				respCode: http.StatusBadRequest,
				resBody:  `{"error":"invalid request body"}`,
			},
		},
		{
			name: "invalid values in request body",
			args: args{
				reqBody: `{"user_id":"u-1"}`,
				ust:     userStorageMock{},
			},
			exp: exp{
				respCode: http.StatusUnprocessableEntity,
				resBody:  `{"errors":[{"path":"user_id","message":"invalid UUID length: 3"}]}`,
			},
		},
		{
			name: "user not found error",
			args: args{
				reqBody: `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe"}`,
				ust: userStorageMock{
					deleteUser: func() error {
						return storage.UserNotFoundErr
					},
				},
			},
			exp: exp{
				respCode: http.StatusNotFound,
				resBody:  `{"error":"user not found"}`,
			},
		},
		{
			name: "database error",
			args: args{
				reqBody: `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe"}`,
				ust: userStorageMock{
					deleteUser: func() error {
						return errors.New("database error")
					},
				},
			},
			exp: exp{
				respCode: http.StatusInternalServerError,
				resBody:  ``,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("", "", strings.NewReader(tc.args.reqBody))
			require.NoError(t, err)

			rec := httptest.NewRecorder()

			h := NewHandler(tc.args.ust)

			h.HandleDeleteUser(rec, req)

			assert.Equal(t, tc.exp.respCode, rec.Code, "unexpected response code")
			assert.Equal(t, tc.exp.resBody, rec.Body.String(), "unexpected response body")
		})
	}
}

func TestParseRequestBody(t *testing.T) {
	type args struct {
		reqBody string
	}
	type exp struct {
		result   bool
		body     configuration.UserIdentifierRequest
		respBody string
		code     int
	}
	tcs := []struct {
		name string
		args args
		exp  exp
	}{
		{
			name: "ok",
			args: args{
				reqBody: `{"user_id":"u-1"}`,
			},
			exp: exp{
				result: true,
				body: configuration.UserIdentifierRequest{
					UserId: "u-1",
				},
				respBody: ``,
				code:     http.StatusOK,
			},
		},
		{
			name: "empty request body",
			args: args{
				reqBody: ``,
			},
			exp: exp{
				result:   false,
				body:     configuration.UserIdentifierRequest{},
				respBody: `{"error":"empty request body"}`,
				code:     http.StatusBadRequest,
			},
		},
		{
			name: "invalid request body",
			args: args{
				reqBody: `"{"id":"xxx"}`,
			},
			exp: exp{
				result:   false,
				body:     configuration.UserIdentifierRequest{},
				respBody: `{"error":"invalid request body"}`,
				code:     http.StatusBadRequest,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var rb configuration.UserIdentifierRequest
			req, err := http.NewRequest("", "", strings.NewReader(tc.args.reqBody))
			require.NoError(t, err)

			rec := httptest.NewRecorder()

			assert.Equal(t, tc.exp.result, ParseRequestBody(rec, req, &rb))
			assert.Equal(t, tc.exp.code, rec.Code)
			assert.Equal(t, tc.exp.respBody, rec.Body.String())
			assert.Equal(t, tc.exp.body, rb)
		})
	}
}

func TestWriteJson(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		expCode := 499

		res := httptest.NewRecorder()
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			WriteJson(
				expCode,
				struct {
					Status string `json:"status"`
				}{Status: "error"},
				w,
			)
		})

		req, err := http.NewRequest(http.MethodPost, "", strings.NewReader(""))
		require.NoError(t, err)

		handler.ServeHTTP(res, req)

		assert.Equal(t, expCode, res.Code, "unexpected code")
		assert.Equal(t, HeaderContentTypeJson, res.Header().Get(HeaderContentType), "unexpected content type")
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
		assert.Equal(t, "", res.Header().Get(HeaderContentType), "unexpected content type")
		assert.Equal(t, "", res.Body.String(), "unexpected body")
	})
}

type userStorageMock struct {
	users      func() ([]storage.User, error)
	user       func() (storage.User, error)
	createUser func() error
	updateUser func() error
	deleteUser func() error
}

func (u userStorageMock) Users() ([]storage.User, error) {
	return u.users()
}

func (u userStorageMock) User(_ string) (storage.User, error) {
	return u.user()
}

func (u userStorageMock) CreateUser(_ storage.User) error {
	return u.createUser()
}

func (u userStorageMock) UpdateUser(_ storage.User) error {
	return u.updateUser()
}

func (u userStorageMock) DeleteUser(_ string) error {
	return u.deleteUser()
}
