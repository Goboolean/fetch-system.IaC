package connect_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.IaC/internal/connect"
	"github.com/Goboolean/fetch-system.IaC/internal/util"
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


func TestSingleTopicConnector(t *testing.T) {

	c := SetupConnect()
	defer TeardownConnect(c)

	const topic = "test.singletopic.connect"

	t.Run("CreateConnector", func(t *testing.T) {
		ctx := context.Background()

		err := c.CreateSingleTopicConnector(ctx, topic)
		assert.NoError(t, err)

		err = c.CheckPluginConfig(ctx, topic)
		assert.NoError(t, err)

		list, err := c.GetConnectors(ctx)
		assert.NoError(t, err)
		assert.Contains(t, list, topic)

		time.Sleep(1 * time.Second)

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


func TestBulkTopicConnector(t *testing.T) {

	c := SetupConnect()
	defer TeardownConnect(c)

	var topics = []string{util.RandomString(10), util.RandomString(10), util.RandomString(10)}
	const name = "test.bulktopic.connect"
	
	t.Run("CreateConnector", func(t *testing.T) {
		ctx := context.Background()

		err := c.CreateBulkTopicConnector(ctx, name, topics)
		assert.NoError(t, err)

		list, err := c.GetConnectors(ctx)
		assert.NoError(t, err)
		assert.Contains(t, list, name)

		err = c.CheckTasksStatus(ctx, name)
		assert.NoError(t, err)
	})

	t.Run("DeleteConnector", func(t *testing.T) {
		ctx := context.Background()

		err := c.DeleteConnector(ctx, name)
		assert.NoError(t, err)

		list, err := c.GetConnectors(ctx)
		assert.NoError(t, err)
		assert.NotContains(t, list, name)
	})
}





func TestConnectorScenario(t *testing.T) {

	// topic containing dots is possible: demonstrated by this test.

	var (
		productId = "test.sibujo.connect"
		timeFrame = "1s"
		topic = fmt.Sprintf("%s.%s", productId, timeFrame)
		times = 10
	)
	
	c := SetupConnect()
	p := SetupProducer()
	a := SetupAdminClient()
	m := SetupMongoClient()

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := c.DeleteAllConnectors(ctx)
		assert.NoError(t, err)

		err = a.DeleteAllTopics(ctx)
		assert.NoError(t, err)

		TeardownConnect(c)
		TeardownProducer(p)
		TeardownAdminClient(a)
	})

	t.Run("CreateConnector", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := c.CreateSingleTopicConnector(ctx, topic)
		assert.NoError(t, err)

		err = c.CheckTasksStatus(ctx, topic)
		assert.NoError(t, err)
	})

	t.Run("ProduceJsonData", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		for i := 0; i < times; i++ {
			var aggregate = connect.Aggregate{
				Timestamp: time.Now().UnixNano() / int64(time.Millisecond),
			}

			payload, err := json.Marshal(aggregate)
			assert.NoError(t, err)

			err = p.Produce(topic, payload)
			assert.NoError(t, err)
		}

		number, err := p.Flush(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, number)
	})

	t.Run("QueryJsonData", func(t *testing.T) {
		time.Sleep(3 * time.Second)

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		data, err := m.FetchAllStockBatch(ctx, productId, timeFrame)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(data), times)
	})
}



func TestBulkTopicConnectorLoad(t *testing.T) {

	c := SetupConnect()
	defer TeardownConnect(c)

	const n = 100
	const name = "test.bulktopicload.connect"

	var topics = make([]string, n)
	for i := 0; i < n; i++ {
		topics[i] = util.RandomString(10)
	}

	t.Run("CheckConnectorNotExists", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		list, err := c.GetConnectors(ctx)
		assert.NoError(t, err)
		assert.NotContains(t, list, name)
	})
	
	t.Run("CreateConnector", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		err := c.CreateBulkTopicConnector(ctx, name, topics)
		assert.NoError(t, err)

		list, err := c.GetConnectors(ctx)
		assert.NoError(t, err)
		assert.Contains(t, list, name)

		err = c.CheckTasksStatus(ctx, name)
		assert.NoError(t, err)
	})

	t.Run("DeleteConnector", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := c.DeleteConnector(ctx, name)
		assert.NoError(t, err)

		list, err := c.GetConnectors(ctx)
		assert.NoError(t, err)
		assert.NotContains(t, list, name)
	})

	t.Run("Sleep to check cpu usage when connector is idle", func(t *testing.T) {
		time.Sleep(10 * time.Second)
	})
}
