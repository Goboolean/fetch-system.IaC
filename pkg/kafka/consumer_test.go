package kafka

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.IaC/pkg/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)



func TestConsumer(t *testing.T) {

	c := SetupConsumer()
	defer TeardownConsumer(c)

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

	const productId = "test.consumeaggs.io"
	const productType = "1s"

	t.Cleanup(func() {
		TeardownConsumer(c)
		TeardownProducer(p)
	})

	const count = 10
	var received = 0

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
			topic := fmt.Sprintf("%s.%s", productId, productType)
			payload, err := proto.Marshal(&model.Aggregate{
				Timestamp: time.Now().UnixNano(),
			})
			assert.NoError(t, err)

			err = p.Produce(topic, payload)
			assert.NoError(t, err)
		}

		_, err := p.Flush(ctx)
		assert.NoError(t, err)
	})

	t.Run("Consume", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		for i := 0; i < count; i++ {
			select {
			case <-ctx.Done():
				assert.Fail(t, "timeout")
				return
			case <-ch:
				received++
			}
		}

		assert.Equal(t, count, received)
	})
}

func TestConsumeTrade(t *testing.T) {

	c := SetupConsumer()
	p := SetupProducer()

	const productId = "test.consumetrade.io"
	const productType = "t"

	t.Cleanup(func() {
		TeardownConsumer(c)
		TeardownProducer(p)
	})

	const count = 10
	var received = 0

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
			topic := fmt.Sprintf("%s.%s", productId, productType)
			payload, err := proto.Marshal(&model.Trade{
				Timestamp: time.Now().UnixNano(),
			})
			assert.NoError(t, err)

			err = p.Produce(topic, payload)
			assert.NoError(t, err)
		}

		_, err := p.Flush(ctx)
		assert.NoError(t, err)
	})

	t.Run("Consume", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		for i := 0; i < count; i++ {
			select {
			case <-ctx.Done():
				assert.Fail(t, "timeout")
				return
			case <-ch:
				received++
			}
		}

		assert.Equal(t, count, received)
	})
}
