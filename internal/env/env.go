package env

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

type InvalidValueError struct {
	key   string
	value string
	err   error
}

func NewInvalidValueError(key string, value string, err error) *InvalidValueError {
	return &InvalidValueError{
		key:   key,
		value: value,
		err:   err,
	}
}

func (e InvalidValueError) Error() string {
	return fmt.Sprintf("invalid env \"%s\" value \"%s\": %s", e.key, e.value, e.err)
}

func Port(key string, required bool, defaultValue string) (int, error) {
	s, err := String(key, required, defaultValue)
	if err != nil {
		return 0, err
	}

	v, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return 0, NewInvalidValueError(key, s, errors.New("port must by in range 0 - 65535"))
	}

	return int(v), nil
}

func Int(key string, required bool, defaultValue string) (int, error) {
	s, err := String(key, required, defaultValue)
	if err != nil {
		return 0, err
	}

	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, NewInvalidValueError(key, s, err)
	}

	return v, nil
}

func MustInt(v int, err error) int {
	if err != nil {
		panic(err)
	}

	return v
}

func String(key string, required bool, defaultValue string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok {
		v = defaultValue
	}

	if required && v == "" {
		return v, NewInvalidValueError(key, v, errors.New("required env variable is not set"))
	}

	return v, nil
}

func MustString(v string, err error) string {
	if err != nil {
		panic(err)
	}

	return v
}
