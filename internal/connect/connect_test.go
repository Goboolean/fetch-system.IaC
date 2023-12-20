package connect_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.IaC/internal/connect"
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

	t.Cleanup(func() {
		err := conf.DeleteAllTopics(context.Background())
		assert.NoError(t, err)

		TeardownConnect(c)
	})

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
		t.Log(err)
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


func TestCreateConnector(t *testing.T) {

	var (
		productId = "test.sibujo.connect"
		timeFrame = "1s"
		topic = fmt.Sprintf("%s.%s", productId, timeFrame)
	)
	
	c := SetupConnect()
	p := SetupProducer()
	a := SetupAdminClient()
	m := SetupMongoClient()

	t.Cleanup(func() {
		ctx := context.Background()

		err := c.DeleteConnector(ctx, topic)
		assert.NoError(t, err)

		err = a.DeleteAllTopics(ctx)
		assert.NoError(t, err)

		TeardownConnect(c)
		TeardownProducer(p)
		TeardownAdminClient(a)
	})

	t.Run("CreateConnector", func(t *testing.T) {
		ctx := context.Background()

		err := c.CreateConnector(ctx, topic)
		assert.NoError(t, err)

		time.Sleep(1 * time.Second)

		err = c.CheckTasksStatus(ctx, topic)
		assert.NoError(t, err)
	})

	t.Run("ProduceJsonData", func(t *testing.T) {

		ctx := context.Background()

		var aggregate = connect.Aggregate{
			Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
		}

		payload, err := json.Marshal(aggregate)
		assert.NoError(t, err)

		err = p.Produce(topic, payload)
		assert.NoError(t, err)

		number, err := p.Flush(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, number)
	})

	t.Run("QueryJsonData", func(t *testing.T) {
		time.Sleep(3 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		data, err := m.FetchAllStockBatch(ctx, productId, timeFrame)
		assert.NoError(t, err)
		assert.NotEmpty(t, data)
	})
}