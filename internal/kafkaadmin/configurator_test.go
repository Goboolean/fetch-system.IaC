package kafkaadmin_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/kafkaadmin"
	"github.com/stretchr/testify/assert"

	_ "github.com/Goboolean/common/pkg/env"
)

func SetupConfigurator() *kafkaadmin.Configurator {

	conf, err := kafkaadmin.New(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
	})
	if err != nil {
		panic(err)
	}

	return conf
}

func TeardownConfigurator(conf *kafkaadmin.Configurator) {
	conf.Close()
}

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

	const topic = "test.createdeletetopic.io.1s"

	conf := SetupConfigurator()
	defer TeardownConfigurator(conf)

	t.Run("CreateTopic", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.CreateTopic(ctx, topic)
		assert.NoError(t, err)

		exists, err := conf.TopicExists(ctx, topic)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("CreateExitingTopic", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.CreateTopic(ctx, topic)
		assert.Error(t, err)
	})

	t.Run("DeleteTopic", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.DeleteTopic(ctx, topic)
		assert.NoError(t, err)

		exists, err := conf.TopicExists(ctx, topic)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("DeleteNonExistingTopic", func(t *testing.T) {

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()

		err := conf.DeleteTopic(ctx, topic)
		assert.Error(t, err)
	})
}


func Test_GetTopicList(t *testing.T) {

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
