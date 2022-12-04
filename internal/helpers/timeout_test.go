package helpers

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_WithTimeout(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		assert.True(t, WithTimeout(func() {
			time.Sleep(time.Second)
		}, time.Second*2))
	})

	t.Run("timeout", func(t *testing.T) {
		assert.False(t, WithTimeout(func() {
			time.Sleep(time.Second * 2)
		}, time.Second))
	})
}
