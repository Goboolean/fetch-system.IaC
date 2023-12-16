package kafka

import (
	"context"
	"sync"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Producer struct {
	producer *kafka.Producer

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

func NewProducer(c *resolver.ConfigMap) (*Producer, error) {

	bootstrap_host, err := c.GetStringKey("BOOTSTRAP_HOST")
	if err != nil {
		return nil, err
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers":   bootstrap_host,
		"acks":                -1,
		"go.delivery.reports": true,
	})

	ctx, cancel := context.WithCancel(context.Background())

	instance := &Producer{
		producer: p,
		wg:       sync.WaitGroup{},
		ctx:      ctx,
		cancel:   cancel,
	}

	return instance, nil
}

func (p *Producer) Produce(topic string, msg []byte) error {

	if err := p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Value:          msg,
	}, nil); err != nil {
		return err
	}

	return nil
}

func (p *Producer) Flush(ctx context.Context) (int, error) {

	go func() {
		for range p.producer.Events() {
		}
	}()

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(1 << 32)
	}

	left := p.producer.Flush(int(time.Until(deadline).Milliseconds()))
	if left != 0 {
		return left, ErrFailedToFlush
	}

	return 0, nil
}

func (p *Producer) Close() {
	p.producer.Close()
	p.cancel()
	p.wg.Wait()
}

func (p *Producer) Ping(ctx context.Context) error {

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(1 << 32)
	}

	remaining := time.Until(deadline)
	_, err := p.producer.GetMetadata(nil, true, int(remaining.Milliseconds()))
	return err
}
