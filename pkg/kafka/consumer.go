package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/reflect/protoreflect"
)



type Consumer struct {
	consumer *kafka.Consumer

	topic string

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// example:
// p, err := NewConsumer[*model.Event](&resolver.ConfigMap{
//   "BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
//   "REGISTRY_HOST":  os.Getenv("KAFKA_REGISTRY_HOST"), // optional
//   "GROUP_ID":       "GROUP_ID",
//   "PROCESSOR_COUNT": os.Getenv("KAFKA_PROCESSOR_COUNT"),
//   "TOPIC":          "TOPIC",
// }, subscribeListenerImpl)
func NewConsumer(c *resolver.ConfigMap) (*Consumer, error) {

	bootstrap_host, err := c.GetStringKey("BOOTSTRAP_HOST")
	if err != nil {
		return nil, err
	}

	group_id, err := c.GetStringKey("GROUP_ID")
	if err != nil {
		return nil, err
	}

	processor_count, err := c.GetIntKey("PROCESSOR_COUNT")
	if err != nil {
		return nil, err
	}

	conn, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":   bootstrap_host,
		"group.id":            group_id,
		"auto.offset.reset": "earliest",
	})

	ctx, cancel := context.WithCancel(context.Background())

	instance := &Consumer{
		consumer: conn,
		wg: sync.WaitGroup{},
		ctx: ctx,
		cancel: cancel,
	}

	go instance.readMessage(ctx, &instance.wg)
	for i := 0; i < processor_count; i++ {
		go instance.consumeMessage(ctx, &instance.wg)
	}
	return instance, nil
}


func (c *Consumer) Subscribe(topic string, schema protoreflect.MessageType) error {
	if c.topic != "" {
		return ErrTopicAlreadySubscribed
	}

	if err := c.consumer.Subscribe(topic, nil); err != nil {
		return err
	}
	c.topic = topic
	return nil
}


func (c *Consumer) readMessage(ctx context.Context, wg *sync.WaitGroup) {
	go func() {
		c.wg.Add(1)
		defer c.wg.Done()

		for {
			if err := c.ctx.Err(); err != nil {
				return
			}

			msg, err := c.consumer.ReadMessage(time.Second * 1)
			if err != nil {
				continue
			}

			var event T
			if err := c.deserial.DeserializeInto(c.topic, msg.Value, event); err != nil {
				log.WithFields(log.Fields{
					"topic": *msg.TopicPartition.Topic,
					"data":  msg.Value,
					"error": err,
				}).Error("Failed to deserialize data")
				continue
			}

			log.WithFields(log.Fields{
				"topic": *msg.TopicPartition.Topic,
				"data":  msg.Value,
				"partition":  msg.TopicPartition.Partition,
				"offset": msg.TopicPartition.Offset,
			}).Trace("Consumer received message")

			c.channel <- event
		}
	}()
}


func (c *Consumer) consumeMessage(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case event := <-c.channel:
			ctx, cancel := context.WithTimeout(c.ctx, time.Second*5)
			if err := c.listener.OnReceiveMessage(ctx, event); err != nil {
				log.WithFields(log.Fields{
					"event":  event,
					"error": err,
				}).Error("Failed to process data")
			}
			cancel()
		}
	}
}


func (c *Consumer) Close() {
	c.cancel()
	time.Sleep(time.Second * 1)
	c.consumer.Close()
	c.wg.Wait()
}


func (c *Consumer) Ping(ctx context.Context) error {
	// It requires ctx to be deadline set, otherwise it will return error
	// It will return error if there is no response within deadline
	deadline, ok := ctx.Deadline()
	if !ok {
		return ErrDeadlineSettingRequired
	}

	remaining := time.Until(deadline)
	_, err := c.consumer.GetMetadata(nil, true, int(remaining.Milliseconds()))
	return err
}