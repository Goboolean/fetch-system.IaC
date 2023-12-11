package connect_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)



func TestConstuctor(t *testing.T) {
	
	c := SetupConnect()

	t.Run("Ping", func(t *testing.T) {
		err := c.Ping(context.Background())
		assert.NoError(t, err)
	})

	t.Run("CheckCompatibility", func(t *testing.T) {
		err := c.CheckCompatibility(context.Background())
		assert.NoError(t, err)
	})
}


func TestConnector(t *testing.T) {

	c := SetupConnect()
	defer TeardownConnect(c)

	const topic = "test.connector.connect"

	t.Run("CreateConnector", func(t *testing.T) {
		ctx := context.Background()

		err := c.CreateConnector(ctx, topic)
		assert.NoError(t, err)

		err = c.CheckPluginConfig(ctx, topic)
		assert.NoError(t, err)

		list, err := c.GetConnectors(ctx)
		assert.NoError(t, err)
		assert.Contains(t, list, topic)

		err = c.CheckTasksStatus(ctx, topic)
		assert.NoError(t, err)
	})

	t.Run("DeleteConnector", func(t *testing.T) {
		ctx := context.Background()

		err := c.DeleteConnector(ctx, topic)
		assert.NoError(t, err)

		list, err := c.GetConnectors(ctx)
		assert.NoError(t, err)
		assert.NotContains(t, list, topic)
	})
}