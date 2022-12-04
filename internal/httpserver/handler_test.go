package httpserver

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"app/internal/configuration"
	"app/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func Test_newHandler(t *testing.T) {
	expHandler := &handler{
		ust: userStorageMock{},
	}

	assert.Equal(t, expHandler, newHandler(userStorageMock{}))
}

func TestHandler_Users(t *testing.T) {
	type args struct {
		ust storage.UserStorage
	}
	type exp struct {
		respCode int
		respBody string
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
				respBody: `{"users":[{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe","name":"John Doe","age":42},{"user_id":"63df08d2-fa53-4575-a681-99058f8daba5","name":"Josh Brave","age":20}]}`,
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
				respBody: ``,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("", "", nil)
			require.NoError(t, err)

			rec := httptest.NewRecorder()

			h := newHandler(tc.args.ust)

			h.Users(rec, req)

			assert.Equal(t, tc.exp.respCode, rec.Code, "unexpected response code")
			assert.Equal(t, tc.exp.respBody, rec.Body.String(), "unexpected response body")
		})
	}
}

func TestHandler_User(t *testing.T) {
	type args struct {
		reqBody string
		ust     storage.UserStorage
	}
	type exp struct {
		respCode int
		respBody string
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
				respBody: `{"user_id":"7df661d5-47e3-4533-baa6-5f952d18bffe","name":"John Doe","age":42}`,
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
				respBody: `{"error":"invalid request body"}`,
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
				respBody: `{"errors":[{"path":"user_id","message":"invalid UUID length: 3"}]}`,
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
				respBody: `{"error":"user not found"}`,
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
				respBody: ``,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("", "", strings.NewReader(tc.args.reqBody))
			require.NoError(t, err)

			rec := httptest.NewRecorder()

			h := newHandler(tc.args.ust)

			h.User(rec, req)

			assert.Equal(t, tc.exp.respCode, rec.Code, "unexpected response code")
			assert.Equal(t, tc.exp.respBody, rec.Body.String(), "unexpected response body")
		})
	}
}

func TestHandler_CreateUser(t *testing.T) {
	type args struct {
		reqBody string
		ust     storage.UserStorage
	}
	type exp struct {
		respCode int
		respBody string
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
				respBody: ``,
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
				respBody: `{"error":"invalid request body"}`,
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
				respBody: `{"errors":[{"path":"user_id","message":"invalid UUID length: 3"}]}`,
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
				respBody: `{"error":"user already exists"}`,
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
				respBody: ``,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("", "", strings.NewReader(tc.args.reqBody))
			require.NoError(t, err)

			rec := httptest.NewRecorder()

			h := newHandler(tc.args.ust)

			h.CreateUser(rec, req)

			assert.Equal(t, tc.exp.respCode, rec.Code, "unexpected response code")
			assert.Equal(t, tc.exp.respBody, rec.Body.String(), "unexpected response body")
		})
	}
}

func TestHandler_UpdateUser(t *testing.T) {
	type args struct {
		reqBody string
		ust     storage.UserStorage
	}
	type exp struct {
		respCode int
		respBody string
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
				respBody: ``,
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
				respBody: `{"error":"invalid request body"}`,
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
				respBody: `{"errors":[{"path":"user_id","message":"invalid UUID length: 3"}]}`,
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
				respBody: `{"error":"user not found"}`,
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
				respBody: ``,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("", "", strings.NewReader(tc.args.reqBody))
			require.NoError(t, err)

			rec := httptest.NewRecorder()

			h := newHandler(tc.args.ust)

			h.UpdateUser(rec, req)

			assert.Equal(t, tc.exp.respCode, rec.Code, "unexpected response code")
			assert.Equal(t, tc.exp.respBody, rec.Body.String(), "unexpected response body")
		})
	}
}

func TestHandler_DeleteUser(t *testing.T) {
	type args struct {
		reqBody string
		ust     storage.UserStorage
	}
	type exp struct {
		respCode int
		respBody string
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
				respBody: ``,
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
				respBody: `{"error":"invalid request body"}`,
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
				respBody: `{"errors":[{"path":"user_id","message":"invalid UUID length: 3"}]}`,
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
				respBody: `{"error":"user not found"}`,
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
				respBody: ``,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("", "", strings.NewReader(tc.args.reqBody))
			require.NoError(t, err)

			rec := httptest.NewRecorder()

			h := newHandler(tc.args.ust)

			h.DeleteUser(rec, req)

			assert.Equal(t, tc.exp.respCode, rec.Code, "unexpected response code")
			assert.Equal(t, tc.exp.respBody, rec.Body.String(), "unexpected response body")
		})
	}
}

func Test_parseRequestBody(t *testing.T) {
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

			assert.Equal(t, tc.exp.result, parseRequestBody(rec, req, &rb))
			assert.Equal(t, tc.exp.code, rec.Code)
			assert.Equal(t, tc.exp.respBody, rec.Body.String())
			assert.Equal(t, tc.exp.body, rb)
		})
	}
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
