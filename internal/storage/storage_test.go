package storage

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

var (
	databaseError = errors.New("database error")

	user1 = User{
		UserId: "7df661d5-47e3-4533-baa6-5f952d18bffe",
		Name:   "John Doe",
		Age:    42,
	}

	user2 = User{
		UserId: "63df08d2-fa53-4575-a681-99058f8daba5",
		Name:   "Josh Brave",
		Age:    20,
	}
)

func TestNewUserStorage(t *testing.T) {
	db, _, _ := sqlmock.New()
	expStorage := &storage{
		db: db,
	}

	assert.Equal(t, expStorage, NewUserStorage(db))
}

func TestStorage_Users(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "name", "age"}).AddRow(user1.UserId, user1.Name, user1.Age).AddRow(user2.UserId, user2.Name, user2.Age)
		mock.ExpectQuery(regexp.QuoteMeta(selectUsersSQL)).WillReturnRows(rows)

		s := NewUserStorage(db)

		users, err := s.Users()

		assert.NoError(t, err)
		assert.Equal(t, []User{user1, user2}, users)
	})

	t.Run("scan error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(user1.UserId, user1.Name)
		mock.ExpectQuery(regexp.QuoteMeta(selectUsersSQL)).WillReturnRows(rows)

		s := NewUserStorage(db)

		res, err := s.Users()

		assert.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(selectUsersSQL)).WillReturnError(databaseError)
		s := NewUserStorage(db)

		res, err := s.Users()

		assert.Equal(t, databaseError, err)
		assert.Nil(t, res)
	})
}

func TestStorage_User(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		rows := sqlmock.NewRows([]string{"id", "name", "age"}).AddRow(user1.UserId, user1.Name, user1.Age)
		mock.ExpectQuery(regexp.QuoteMeta(selectUserSQL)).WithArgs(user1.UserId).WillReturnRows(rows)

		s := NewUserStorage(db)

		user, err := s.User(user1.UserId)

		assert.NoError(t, err)
		assert.Equal(t, user1, user)
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(selectUserSQL)).WithArgs(user1.UserId).WillReturnError(databaseError)

		s := NewUserStorage(db)

		_, err = s.User(user1.UserId)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("not found error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectQuery(regexp.QuoteMeta(selectUserSQL)).WithArgs(user1.UserId).WillReturnError(sql.ErrNoRows)

		s := NewUserStorage(db)

		_, err = s.User(user1.UserId)

		assert.Equal(t, UserNotFoundErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStorage_CreateUser(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(insertUserSQL)).WithArgs(user1.UserId, user1.Name, user1.Age).WillReturnResult(sqlmock.NewResult(0, 1))

		s := NewUserStorage(db)

		require.NoError(t, s.CreateUser(user1))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("entry already exists", func(t *testing.T) {
		expErr := errors.New(`pq: duplicate key value violates unique constraint "user_id_pk"`)

		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(insertUserSQL)).WithArgs(user1.UserId, user1.Name, user1.Age).WillReturnError(expErr)

		s := NewUserStorage(db)

		err = s.CreateUser(user1)

		assert.Equal(t, UserAlreadyExistsErr, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(insertUserSQL)).WithArgs(user1.UserId, user1.Name, user1.Age).WillReturnError(databaseError)

		s := NewUserStorage(db)

		err = s.CreateUser(user1)

		assert.Equal(t, databaseError, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStorage_UpdateUser(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(updateUserSQL)).WithArgs(user1.UserId, user1.Name, user1.Age).WillReturnResult(sqlmock.NewResult(0, 1))

		s := NewUserStorage(db)

		require.NoError(t, s.UpdateUser(user1))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(updateUserSQL)).WithArgs(user1.UserId, user1.Name, user1.Age).WillReturnResult(sqlmock.NewResult(0, 0))

		s := NewUserStorage(db)

		assert.Equal(t, UserNotFoundErr, s.UpdateUser(user1))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows affected error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(updateUserSQL)).WithArgs(user1.UserId, user1.Name, user1.Age).WillReturnResult(sqlmock.NewErrorResult(databaseError))

		s := NewUserStorage(db)

		assert.Equal(t, databaseError, s.UpdateUser(user1))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(updateUserSQL)).WillReturnError(databaseError)

		s := NewUserStorage(db)

		assert.Equal(t, databaseError, s.UpdateUser(user1))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestStorage_DeleteUser(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(deleteUserSQL)).WithArgs(user1.UserId).WillReturnResult(sqlmock.NewResult(0, 1))

		s := NewUserStorage(db)

		require.NoError(t, s.DeleteUser(user1.UserId))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("user not found", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(deleteUserSQL)).WithArgs(user1.UserId).WillReturnResult(sqlmock.NewResult(0, 0))

		s := NewUserStorage(db)

		assert.Equal(t, UserNotFoundErr, s.DeleteUser(user1.UserId))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("rows affected error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(deleteUserSQL)).WithArgs(user1.UserId).WillReturnResult(sqlmock.NewErrorResult(databaseError))

		s := NewUserStorage(db)

		assert.Equal(t, databaseError, s.DeleteUser(user1.UserId))
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("db error", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		mock.ExpectExec(regexp.QuoteMeta(deleteUserSQL)).WillReturnError(databaseError)

		s := NewUserStorage(db)

		assert.Equal(t, databaseError, s.DeleteUser(user1.UserId))
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestAlreadyExistsErr(t *testing.T) {
	t.Run("already exists error", func(t *testing.T) {
		assert.True(t, AlreadyExistsErr(errors.New("pq: duplicate key value violates unique constraint")))
	})

	t.Run("not already exists error", func(t *testing.T) {
		assert.False(t, AlreadyExistsErr(errors.New("pq: duplicate key value violates unique")))
	})
}
