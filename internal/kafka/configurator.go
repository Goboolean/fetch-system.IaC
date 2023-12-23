package kafka

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/util"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	log "github.com/sirupsen/logrus"
)

// Configurator has a role for making and deleting topic, checking topic exists, and getting topic list.
type Configurator struct {
	client *kafka.AdminClient

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

var (
	conf *Configurator
	once sync.Once
)

func New(c *resolver.ConfigMap) (*Configurator, error) {

	bootstrap_host, err := c.GetStringKey("BOOTSTRAP_HOST")
	if err != nil {
		return nil, err
	}

	config := &kafka.ConfigMap{
		"bootstrap.servers": bootstrap_host,
	}

	conn, err := kafka.NewAdminClient(config)

	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Configurator{
		client: conn,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (c *Configurator) Close() {
	c.cancel()
	time.Sleep(time.Second * 1)
	c.client.Close()
	c.wg.Wait()
}

func (c *Configurator) ping(ctx context.Context) error {
	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(time.Hour)
	}

	_, err := c.client.GetMetadata(nil, true, int(time.Until(deadline).Milliseconds()))
	return err
}

func (c *Configurator) Ping(ctx context.Context) error {
	for {
		if err := c.ping(ctx); err != nil {
			log.WithField("error", err).Error("Failed to ping, waiting 5 seconds")
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(5 * time.Second):
				continue
			}
		}

		return nil
	}
}

func (c *Configurator) CreateTopic(ctx context.Context, topic string) error {

	topicInfo := kafka.TopicSpecification{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}

	result, err := c.client.CreateTopics(ctx, []kafka.TopicSpecification{topicInfo})
	if err != nil {
		return err
	}

	if err := result[0].Error; err.Code() != kafka.ErrNoError {
		if err.Code() == kafka.ErrTopicAlreadyExists {
			return nil
		}
		return fmt.Errorf(err.String())
	}
	return nil
}


func (c *Configurator) CreateTopics(ctx context.Context, topics ...string) error {

	topicInfos := make([]kafka.TopicSpecification, len(topics))
	for i, topic := range topics {
		topicInfos[i] = kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		}
	}

	result, err := c.client.CreateTopics(ctx, topicInfos)
	if err != nil {
		return err
	}

	for _, r := range result {
		if err := r.Error; err.Code() != kafka.ErrNoError {
			if err.Code() == kafka.ErrTopicAlreadyExists {
				continue
			}
			return fmt.Errorf(err.String())
		}
	}
	return nil
}



func (c *Configurator) DeleteTopic(ctx context.Context, topic string) error {

	result, err := c.client.DeleteTopics(ctx, []string{topic})
	if err != nil {
		return err
	}

	if err := result[0].Error; err.Code() != kafka.ErrNoError {
		return err
	}

	return nil
}

func (c *Configurator) DeleteAllTopics(ctx context.Context) error {
	topicList, err := c.GetTopicList(ctx)
	if err != nil {
		return err
	}

	filteredTopicList := make([]string, 0)
	for _, topic := range topicList {
		if util.Contains(defaultTopicList, topic) {
			continue
		}
		filteredTopicList = append(filteredTopicList, topic)
	}

	if len(filteredTopicList) == 0 {
		return nil
	}

	result, err := c.client.DeleteTopics(ctx, filteredTopicList)
	if err != nil {
		return err
	}

	if err := result[0].Error; err.Code() != kafka.ErrNoError {
		return fmt.Errorf(err.String())
	}

	return nil
}

func (c *Configurator) TopicExists(ctx context.Context, topic string) (bool, error) {

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(1 << 32)
	}

	metadata, err := c.client.GetMetadata(nil, true, int(time.Until(deadline).Milliseconds()))
	if err != nil {
		return false, err
	}

	detail, exists := metadata.Topics[topic]
	detail.Error.Code()
	return exists && detail.Topic != "", nil
}

func (c *Configurator) AllTopicExists(ctx context.Context, topics ...string) (bool, error) {

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(1 << 32)
	}

	metadata, err := c.client.GetMetadata(nil, true, int(time.Until(deadline).Milliseconds()))
	if err != nil {
		return false, err
	}

	for _, topic := range topics {
		detail, exists := metadata.Topics[topic]
		if !exists || detail.Topic == "" {
			return false, nil
		}
	}

	return true, nil
}


func (c *Configurator) GetTopicList(ctx context.Context) ([]string, error) {

	deadline, ok := ctx.Deadline()
	if !ok {
		deadline = time.Now().Add(1 << 32)
	}

	metadata, err := c.client.GetMetadata(nil, true, int(time.Until(deadline).Milliseconds()))
	if err != nil {
		return nil, err
	}

	topicList := make([]string, 0)

	for topic := range metadata.Topics {
		if topic == "" {
			continue
		}
		if util.Contains(defaultTopicList, topic) {
			continue
		}
		topicList = append(topicList, topic)
	}

	return topicList, nil
}

var defaultTopicList = []string{"__consumer_offsets", "kafka-connect-offsets", "kafka-connect-status", "kafka-connect-configs"}
