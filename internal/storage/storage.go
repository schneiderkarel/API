package storage

import (
	"database/sql"
	"errors"
	"strings"
)

const (
	selectUsersSQL = `SELECT user_id, "name", age FROM user_service.users`
	selectUserSQL  = `SELECT user_id, "name", age FROM user_service.users WHERE user_id = $1`
	insertUserSQL  = `INSERT INTO user_service.users (user_id, "name", age) VALUES ($1, $2, $3)`
	updateUserSQL  = `UPDATE user_service.users SET "name" = $2, age = $3 WHERE user_id = $1`
	deleteUserSQL  = `DELETE FROM user_service.users WHERE user_id = $1`
)

var (
	UserAlreadyExistsErr = errors.New("user already exists")
	UserNotFoundErr      = errors.New("user not found")
)

type User struct {
	UserId string `json:"user_id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
}

type UsersResponse struct {
	Users []User `json:"users"`
}

type UserStorage interface {
	Users() ([]User, error)
	User(userId string) (User, error)
	CreateUser(usr User) error
	UpdateUser(usr User) error
	DeleteUser(userId string) error
}

type storage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) UserStorage {
	return &storage{
		db: db,
	}
}

func (st *storage) Users() ([]User, error) {
	var users []User

	rows, err := st.db.Query(selectUsersSQL)
	if err != nil {
		return nil, err
	}

	var usr User
	for rows.Next() {
		if err := rows.Scan(
			&usr.UserId,
			&usr.Name,
			&usr.Age,
		); err != nil {
			return nil, err
		}

		users = append(users, usr)
	}

	return users, err
}

func (st *storage) User(userId string) (User, error) {
	row := st.db.QueryRow(selectUserSQL, userId)

	var usr User
	if err := row.Scan(
		&usr.UserId,
		&usr.Name,
		&usr.Age,
	); err != nil {
		if err == sql.ErrNoRows {
			return User{}, UserNotFoundErr
		}

		return User{}, err
	}

	return usr, nil
}

func (st *storage) CreateUser(usr User) error {
	if _, err := st.db.Exec(insertUserSQL, usr.UserId, usr.Name, usr.Age); err != nil {
		if AlreadyExistsErr(err) {
			return UserAlreadyExistsErr
		}

		return err
	}

	return nil
}

func (st *storage) UpdateUser(usr User) error {
	res, err := st.db.Exec(updateUserSQL, usr.UserId, usr.Name, usr.Age)
	if err != nil {
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return UserNotFoundErr
	}

	return nil
}

func (st *storage) DeleteUser(userId string) error {
	res, err := st.db.Exec(deleteUserSQL, userId)
	if err != nil {
		return err
	}

	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return UserNotFoundErr
	}

	return nil
}

func AlreadyExistsErr(err error) bool {
	return strings.HasPrefix(err.Error(), "pq: duplicate key value violates unique constraint")
}
