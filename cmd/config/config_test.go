package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	os.Clearenv()

	require.NoError(t, os.Setenv(envKeyHttpServerPort, "81"))
	require.NoError(t, os.Setenv(envKeyPostgresDSN, "host=postgres port=5432 user=postgres dbname=postgres sslmode=disable"))

	t.Run("ok", func(t *testing.T) {
		expCfg := Config{
			HttpServerPort: 81,
			PostgresDSN:    "host=postgres port=5432 user=postgres dbname=postgres sslmode=disable",
		}

		require.NotPanics(t, func() {
			assert.Equal(t, expCfg, *Init())
		})
	})

	t.Run("error", func(t *testing.T) {
		require.NoError(t, os.Unsetenv(envKeyPostgresDSN))

		assert.Panics(t, func() {
			Init()
		})
	})
}
