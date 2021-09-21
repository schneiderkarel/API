package env

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewInvalidValueError(t *testing.T) {
	expNewInvalidValueErr := &InvalidValueError{
		key:   "k",
		value: "v",
		err:   errors.New("error"),
	}

	assert.Equal(t, expNewInvalidValueErr, NewInvalidValueError("k", "v", errors.New("error")))
}

func TestInvalidValueError_Error(t *testing.T) {
	invalidValueErr := &InvalidValueError{
		key:   "k",
		value: "v",
		err:   errors.New("error"),
	}

	assert.Equal(t, "invalid env \"k\" value \"v\": error", invalidValueErr.Error())
}

func TestPort(t *testing.T) {
	const (
		envKeyFilled  = "TEST_FILLED"
		envKeyInvalid = "TEST_INVALID"
		envKeyEmpty   = "TEST_EMPTY"
		envKeyNotSet  = "TEST_NOT_SET"
	)

	os.Clearenv()

	require.NoError(t, os.Setenv(envKeyFilled, "65535"))
	require.NoError(t, os.Setenv(envKeyInvalid, "65536"))
	require.NoError(t, os.Setenv(envKeyEmpty, ""))

	type args struct {
		key          string
		required     bool
		defaultValue string
	}
	type exp struct {
		value int
		error bool
	}
	tcs := []struct {
		name string
		args args
		exp  exp
	}{
		{
			name: "filled",
			args: args{
				key: envKeyFilled, required: true, defaultValue: "",
			},
			exp: exp{
				value: 65535, error: false,
			},
		},
		{
			name: "invalid",
			args: args{
				key: envKeyInvalid, required: false, defaultValue: "1",
			},
			exp: exp{
				value: 0, error: true,
			},
		},
		{
			name: "empty with default",
			args: args{
				key: envKeyEmpty, required: true, defaultValue: "1",
			},
			exp: exp{
				value: 0, error: true,
			},
		},
		{
			name: "unset with default",
			args: args{
				key: envKeyNotSet, required: true, defaultValue: "1",
			},
			exp: exp{
				value: 1, error: false,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			v, err := Port(tc.args.key, tc.args.required, tc.args.defaultValue)

			if tc.exp.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.exp.value, v)
		})
	}
}

func TestInt(t *testing.T) {
	const (
		envKeyFilled  = "TEST_FILLED"
		envKeyInvalid = "TEST_INVALID"
		envKeyEmpty   = "TEST_EMPTY"
		envKeyNotSet  = "TEST_NOT_SET"
	)

	os.Clearenv()

	require.NoError(t, os.Setenv(envKeyFilled, "-32"))
	require.NoError(t, os.Setenv(envKeyInvalid, "3.25"))
	require.NoError(t, os.Setenv(envKeyEmpty, ""))

	type args struct {
		key          string
		required     bool
		defaultValue string
	}
	type exp struct {
		value int
		error bool
	}
	tcs := []struct {
		name string
		args args
		exp  exp
	}{
		{
			name: "filled",
			args: args{
				key: envKeyFilled, required: true, defaultValue: "",
			},
			exp: exp{
				value: -32, error: false,
			},
		},
		{
			name: "invalid",
			args: args{
				key: envKeyInvalid, required: false, defaultValue: "551",
			},
			exp: exp{
				value: 0, error: true,
			},
		},
		{
			name: "empty with default",
			args: args{
				key: envKeyEmpty, required: true, defaultValue: "551",
			},
			exp: exp{
				value: 0, error: true,
			},
		},
		{
			name: "unset with default",
			args: args{
				key: envKeyNotSet, required: true, defaultValue: "551",
			},
			exp: exp{
				value: 551, error: false,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			v, err := Int(tc.args.key, tc.args.required, tc.args.defaultValue)

			if tc.exp.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.exp.value, v)
		})
	}
}

func TestMustInt(t *testing.T) {
	os.Clearenv()

	t.Run("ok", func(t *testing.T) {
		assert.NotPanics(t, func() {
			MustInt(Int("TEST_NOT_SET", false, "0"))
		})
	})

	t.Run("error", func(t *testing.T) {
		assert.Panics(t, func() {
			MustInt(Int("TEST_NOT_SET", true, ""))
		})
	})
}

func TestString(t *testing.T) {
	const (
		envKeyFilled = "TEST_FILLED"
		envKeyEmpty  = "TEST_EMPTY"
		envKeyNotSet = "TEST_NOT_SET"
	)

	os.Clearenv()

	require.NoError(t, os.Setenv(envKeyFilled, "tst"))
	require.NoError(t, os.Setenv(envKeyEmpty, ""))

	type args struct {
		key          string
		required     bool
		defaultValue string
	}
	type exp struct {
		value string
		err   bool
	}
	tcs := []struct {
		name string
		args args
		exp  exp
	}{
		{
			name: "filled",
			args: args{
				key: envKeyFilled, required: true, defaultValue: "",
			},
			exp: exp{
				value: "tst", err: false,
			},
		},
		{
			name: "filled with default",
			args: args{
				key: envKeyFilled, required: true, defaultValue: "def",
			},
			exp: exp{
				value: "tst", err: false,
			},
		},
		{
			name: "empty",
			args: args{
				key: envKeyEmpty, required: false, defaultValue: "",
			},
			exp: exp{
				value: "", err: false,
			},
		},
		{
			name: "empty with default",
			args: args{
				key: envKeyEmpty, required: false, defaultValue: "def",
			},
			exp: exp{
				value: "", err: false,
			},
		},
		{
			name: "empty required",
			args: args{
				key: envKeyEmpty, required: true, defaultValue: "",
			},
			exp: exp{
				value: "", err: true,
			},
		},
		{
			name: "empty required with default",
			args: args{
				key: envKeyEmpty, required: true, defaultValue: "def",
			},
			exp: exp{
				value: "", err: true,
			},
		},
		{
			name: "unset",
			args: args{
				key: envKeyNotSet, required: false, defaultValue: "",
			},
			exp: exp{
				value: "", err: false,
			},
		},
		{
			name: "unset with default",
			args: args{
				key: envKeyNotSet, required: false, defaultValue: "def",
			},
			exp: exp{
				value: "def", err: false,
			},
		},
		{
			name: "unset required",
			args: args{
				key: envKeyNotSet, required: true, defaultValue: "",
			},
			exp: exp{
				value: "", err: true,
			},
		},
		{
			name: "unset required with default",
			args: args{
				key: envKeyNotSet, required: true, defaultValue: "def",
			},
			exp: exp{
				value: "def", err: false,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			v, err := String(tc.args.key, tc.args.required, tc.args.defaultValue)

			if tc.exp.err {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tc.exp.value, v)
		})
	}
}

func TestMustString(t *testing.T) {
	os.Clearenv()

	t.Run("ok", func(t *testing.T) {
		assert.NotPanics(t, func() {
			MustString(String("TEST_NOT_SET", false, ""))
		})
	})

	t.Run("error", func(t *testing.T) {
		assert.Panics(t, func() {
			MustString(String("TEST_NOT_SET", true, ""))
		})
	})
}
