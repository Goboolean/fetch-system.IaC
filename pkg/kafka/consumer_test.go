package kafka_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/pkg/model"
	"github.com/Goboolean/fetch-system.IaC/pkg/kafka"
	"github.com/stretchr/testify/assert"
)

func SetupConsumer() *kafka.Consumer {

	c, err := kafka.NewConsumer(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
		"GROUP_ID":       "TEST_GROUP",
	})
	if err != nil {
		panic(err)
	}
	return c
}

func TeardownConsumer(c *kafka.Consumer) {
	mutex.Lock()
	defer mutex.Unlock()
	c.Close()
}

func TestConsumer(t *testing.T) {

	c := SetupConsumer()
	a := SetupConfigurator()

	t.Cleanup(func() {
		err := a.DeleteAllTopics(context.Background())
		assert.NoError(t, err)

		TeardownConsumer(c)
		TeardownConfigurator(a)
	})

	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		err := c.Ping(ctx)
		assert.NoError(t, err)
	})
}

func TestConsumeAggs(t *testing.T) {

	c := SetupConsumer()
	p := SetupProducer()
	a := SetupConfigurator()

	const productId = "test.goboolean.io"
	const productType = "1s"

	t.Cleanup(func() {
		err := a.DeleteAllTopics(context.Background())
		assert.NoError(t, err)

		TeardownConsumer(c)
		TeardownProducer(p)
		TeardownConfigurator(a)
	})

	const count = 3

	var ch <-chan *model.Aggregate

	t.Run("Subscribe", func(t *testing.T) {
		var err error
		ch, err = c.SubscribeAggs(productId, productType)
		assert.NoError(t, err)
		assert.NotNil(t, ch)
	})

	t.Run("Produce", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		for i := 0; i < count; i++ {
			err := p.ProduceAggs(productId, productType, &model.Aggregate{
				Timestamp: time.Now().UnixNano(),
			})
			assert.NoError(t, err)
		}

		_, err := p.Flush(ctx)
		assert.NoError(t, err)
	})

	t.Run("Consume", func(t *testing.T) {
		time.Sleep(time.Second * 5)
		assert.Equal(t, count, len(ch))
	})
}

func TestConsumeTrade(t *testing.T) {

	c := SetupConsumer()
	p := SetupProducer()
	a := SetupConfigurator()

	const productId = "test.goboolean.io"
	const productType = "1m"

	t.Cleanup(func() {
		err := a.DeleteAllTopics(context.Background())
		assert.NoError(t, err)

		TeardownConsumer(c)
		TeardownProducer(p)
		TeardownConfigurator(a)
	})

	const count = 3

	var ch <-chan *model.Trade

	t.Run("Subscribe", func(t *testing.T) {
		var err error
		ch, err = c.SubscribeTrade(productId)
		assert.NoError(t, err)
		assert.NotNil(t, ch)
	})

	t.Run("Produce", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		for i := 0; i < count; i++ {
			err := p.ProduceTrade(productId, &model.Trade{
				Timestamp: time.Now().UnixNano(),
			})
			assert.NoError(t, err)
		}

		_, err := p.Flush(ctx)
		assert.NoError(t, err)
	})

	t.Run("Consume", func(t *testing.T) {
		time.Sleep(time.Second * 5)
		assert.Equal(t, count, len(ch))
	})
}
