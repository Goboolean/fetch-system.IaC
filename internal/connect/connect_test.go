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

	_ "github.com/Goboolean/common/pkg/env"
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

	const (
		topic = "test.singletopic.connect"
		name = "test.singletopic.connect"
		tasks = 3
	)

	var config = connect.ConnectorTopicConfig{
		Topic: topic,
		Collection: topic,
		RotateIntervalMs: 1000,			
	}

	t.Run("CreateConnector", func(t *testing.T) {
		ctx := context.Background()

		err := c.CreateSingleTopicConnector(ctx, name, tasks, config)
		assert.NoError(t, err)

		err = c.CheckPluginConfig(ctx, topic)
		assert.NoError(t, err)

		list, err := c.GetConnectors(ctx)
		assert.NoError(t, err)
		assert.Contains(t, list, topic)

		count, err := c.CheckTasksStatus(ctx, topic)
		assert.NoError(t, err)
		assert.NotZero(t, count)
		t.Log("Number of tasks", count)
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

	const (
		n = 3
		name = "bulktopic"
		tasks = 3
	)

	var configs = make([]connect.ConnectorTopicConfig, n)
	for i := 0; i < n; i++ {
		topic := util.RandomString(10)
		configs[i] = connect.ConnectorTopicConfig{
			Topic: topic,
			Collection: topic,
			RotateIntervalMs: 100000,
		}
	}
	
	t.Run("CreateConnector", func(t *testing.T) {
		ctx := context.Background()

		err := c.CreateBulkTopicConnector(ctx, name, tasks, configs)
		assert.NoError(t, err)

		list, err := c.GetConnectors(ctx)
		assert.NoError(t, err)
		assert.Contains(t, list, name)

		count, err := c.CheckTasksStatus(ctx, name)
		assert.NoError(t, err)
		assert.NotZero(t, count)
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

	const (
		productId = "test.helloworld.connect"
		timeFrame = "t"
		name = "connectorscenario"
		tasks = 3

		times = 10
	)

	var topic = fmt.Sprintf("%s.%s", productId, timeFrame)

	var config = connect.ConnectorTopicConfig{
		Topic: topic,
		Collection: topic,
		RotateIntervalMs: 1000,	
	}
	
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

		err := c.CreateSingleTopicConnector(ctx, name, tasks, config)
		assert.NoError(t, err)

		count, err := c.CheckTasksStatus(ctx, name)
		assert.NoError(t, err)
		assert.NotZero(t, count)
		t.Log("Number of tasks", count)
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
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		for {
			data, err := m.FetchAllStockBatch(ctx, productId, timeFrame)
			assert.NoError(t, err)
			if err != nil {
				return
			}

			if len(data) == times {
				assert.Equal(t, times, len(data))
				break
			}
		}
	})
}



func TestBulkTopicConnectorLoad(t *testing.T) {

	t.Skip("Skip this test because it takes too long time")

	c := SetupConnect()
	defer TeardownConnect(c)

	consumer := SetupConsumer()
	defer TeardownConsumer(consumer)

	const (
		n = 50
		name = "test.bulktopicload.connect"
		tasks = 10
	)

	var configs = make([]connect.ConnectorTopicConfig, n)
	for i := 0; i < n; i++ {
		topic := util.RandomString(10)
		configs[i] = connect.ConnectorTopicConfig{
			Topic: topic,
			Collection: topic,
			RotateIntervalMs: 100000,
		}
	}

	var topics = make([]string, n)
	for i := 0; i < n; i++ {
		topics[i] = util.RandomString(10)
	}

	t.Run("CheckConnectorNotExists", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		list, err := c.GetConnectors(ctx)
		assert.NoError(t, err)
		if util.Contains[string](list, name) {
			err := c.DeleteConnector(ctx, name)
			assert.NoError(t, err)
		}
	})
	
	t.Run("CreateConnector", func(t *testing.T) {
		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		err := c.CreateBulkTopicConnector(ctx, name, tasks, configs)
		assert.NoError(t, err)

		list, err := c.GetConnectors(ctx)
		assert.NoError(t, err)
		assert.Contains(t, list, name)

		count, err := c.CheckTasksStatus(ctx, name)
		assert.NoError(t, err)
		assert.NotZero(t, count)

		elasped := time.Since(start)
		t.Log("Elapsed time:", elasped)
		t.Log("Number of tasks", count)
	})

	t.Run("DeleteConnector", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := c.DeleteConnector(ctx, name)
		assert.NoError(t, err)

		list, err := c.GetConnectors(ctx)
		assert.NoError(t, err)
		assert.NotContains(t, list, name)
	})
}



