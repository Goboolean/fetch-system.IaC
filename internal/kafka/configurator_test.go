package kafka_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)



func Test_Configurator(t *testing.T) {

	conf := SetupConfigurator()
	defer TeardownConfigurator(conf)

	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := conf.Ping(ctx)
		assert.NoError(t, err)
	})
}

func Test_CreateDeleteTopic(t *testing.T) {

	//t.Skip("Skip this test since unhandlable error occurs on kafka admin client.")

	var topic = []string{
		"test.createdeletetopic.io.t",
		"test.createdeletetopic.io.1s",
		"test.createdeletetopic.io.1m",
	}

	conf := SetupConfigurator()
	defer TeardownConfigurator(conf)

	t.Run("CreateTopic", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.CreateTopic(ctx, topic[0])
		assert.NoError(t, err)

		exists, err := conf.TopicExists(ctx, topic[0])
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("CreateExitingTopic", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.CreateTopic(ctx, topic[0])
		assert.NoError(t, err)
	})

	t.Run("CreateMultipleTopic", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.CreateTopics(ctx, topic[1:]...)
		assert.NoError(t, err)

		exists, err := conf.AllTopicExists(ctx, topic...)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("DeleteTopic", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.DeleteTopic(ctx, topic[2])
		assert.NoError(t, err)

		exists, err := conf.TopicExists(ctx, topic[2])
		assert.NoError(t, err)
		assert.False(t, exists)

		exists, err = conf.AllTopicExists(ctx, topic...)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("DeleteNonExistingTopic", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.DeleteTopic(ctx, topic[2])
		assert.Error(t, err)
	})

	t.Run("DeleteAllTopics", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.DeleteAllTopics(ctx)
		assert.NoError(t, err)
	})
}

func Test_GetTopicList(t *testing.T) {

	//t.Skip("Skip this test since unhandlable error occurs on kafka admin client.")

	conf := SetupConfigurator()
	defer TeardownConfigurator(conf)

	topicList := []string{"test.gettopiclist.io.1s", "test.gettopiclist.io.1t"}

	t.Run("CreateTopics", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		for _, topic := range topicList {
			err := conf.CreateTopic(ctx, topic)
			assert.NoError(t, err)
		}
	})

	t.Run("GetTopicList", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		list, err := conf.GetTopicList(ctx)
		assert.NoError(t, err)

		assert.Len(t, list, len(topicList))
		for _, topic := range topicList {
			assert.Contains(t, list, topic)
		}
	})

	t.Run("DeleteAllTopics", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.DeleteAllTopics(ctx)
		assert.NoError(t, err)

		list, err := conf.GetTopicList(ctx)
		assert.NoError(t, err)
		assert.Len(t, list, 0)
	})
}
