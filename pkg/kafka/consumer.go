package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/pkg/model"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

type Consumer struct {
	consumer *kafka.Consumer

	topic string

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

// example:
//
//	p, err := NewConsumer[*model.Event](&resolver.ConfigMap{
//	  "BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
//	  "REGISTRY_HOST":  os.Getenv("KAFKA_REGISTRY_HOST"), // optional
//	  "GROUP_ID":       "GROUP_ID",
//	  "PROCESSOR_COUNT": os.Getenv("KAFKA_PROCESSOR_COUNT"),
//	  "TOPIC":          "TOPIC",
//	}, subscribeListenerImpl)
func NewConsumer(c *resolver.ConfigMap) (*Consumer, error) {

	bootstrap_host, err := c.GetStringKey("BOOTSTRAP_HOST")
	if err != nil {
		return nil, err
	}

	group_id, err := c.GetStringKey("GROUP_ID")
	if err != nil {
		return nil, err
	}

	conn, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrap_host,
		"group.id":          group_id,
		"auto.offset.reset": "earliest",
	})

	ctx, cancel := context.WithCancel(context.Background())

	instance := &Consumer{
		consumer: conn,
		wg:       sync.WaitGroup{},
		ctx:      ctx,
		cancel:   cancel,
	}

	return instance, nil
}

func (c *Consumer) SubscribeTrade(productId string) (<-chan *model.TradeJson, error) {

	topic := fmt.Sprintf("%s.t", productId)
	if !model.IsSymbolValid(topic) {
		return nil, ErrInvalidSymbol
	}

	if err := c.consumer.Subscribe(topic, nil); err != nil {
		return nil, err
	}

	channel := make(chan *model.TradeJson, 100)

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

			var trade model.TradeJson
			if err := json.Unmarshal(msg.Value, &trade); err != nil {
				log.WithFields(log.Fields{
					"topic": *msg.TopicPartition.Topic,
					"data":  msg.Value,
					"error": err,
				}).Error("Failed to deserialize data")
			}

			channel <- &trade
		}
	}()

	return channel, nil
}

func (c *Consumer) SubscribeAggs(productId string, productType string) (<-chan *model.AggregateJson, error) {

	topic := fmt.Sprintf("%s.%s", productId, productType)
	if !model.IsSymbolValid(topic) {
		return nil, ErrInvalidSymbol
	}

	if err := c.consumer.Subscribe(topic, nil); err != nil {
		return nil, err
	}

	channel := make(chan *model.AggregateJson, 100)

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

			var aggs model.AggregateJson
			if err := json.Unmarshal(msg.Value, &aggs); err != nil {
				log.WithFields(log.Fields{
					"topic": *msg.TopicPartition.Topic,
					"data":  msg.Value,
					"error": err,
				}).Error("Failed to deserialize data")
			}

			channel <- &aggs
		}
	}()

	return channel, nil
}

func (c *Consumer) Close() {
	c.cancel()
	time.Sleep(time.Second * 1)
	c.consumer.Close()
	c.wg.Wait()
}

func (c *Consumer) ping(ctx context.Context) error {
	// It requires ctx to be deadline set, otherwise it will return error
	// It will return error if there is no response within deadline
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(1 << 32)
	}

	remaining := time.Until(deadline)
	_, err := c.consumer.GetMetadata(nil, true, int(remaining.Milliseconds()))
	return err
}

func (c *Consumer) Ping(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if err := c.ping(ctx); err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		return nil
	}
}