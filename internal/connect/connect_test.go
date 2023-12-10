package connect_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)



func TestConnect(t *testing.T) {
	
	c := SetupConnect()

	t.Run("Ping", func(t *testing.T) {
		err := c.Ping(context.Background())
		assert.NoError(t, err)
	})
}

