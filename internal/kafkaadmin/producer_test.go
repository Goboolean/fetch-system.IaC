package kafkaadmin_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/kafkaadmin"
	"github.com/stretchr/testify/assert"
)

func SetupProducer() *kafkaadmin.Producer {
	p, err := kafkaadmin.NewProducer(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
	})
	if err != nil {
		panic(err)
	}
	return p
}

func TeardownProducer(p *kafkaadmin.Producer) {
	mutex.Lock()
	defer mutex.Unlock()
	p.Close()
}

func Test_Producer(t *testing.T) {

	p := SetupProducer()
	conf := SetupConfigurator()

	t.Cleanup(func() {
		err := conf.DeleteAllTopics(context.Background())
		assert.NoError(t, err)

		TeardownProducer(p)
	})

	const topic = "test.producer.abc"

	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		err := p.Ping(ctx)
		assert.NoError(t, err)
	})

	t.Run("Produce", func(t *testing.T) {
		const count = 10

		for i := 0; i < count; i++ {
			var data map[string]string = make(map[string]string)
			data["timestamp"] = time.Now().Format(time.RFC3339Nano)
	
			payload, err := json.Marshal(data)
			assert.NoError(t, err)

			err = p.Produce(topic, payload)
			assert.NoError(t, err)
		}
	})

	t.Run("Flush", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*7)
		defer cancel()

		count, err := p.Flush(ctx)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}
