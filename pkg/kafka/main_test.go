package kafka

import (
	"context"
	"os"
	"sync"
	"testing"

	_ "github.com/Goboolean/common/pkg/env"
	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/kafka"
)

var mutex = &sync.Mutex{}




func SetupConfigurator() *kafka.Configurator {
	c, err := kafka.New(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
	})
	if err != nil {
		panic(err)
	}
	return c
}

func TeardownConfigurator(c *kafka.Configurator) {
	mutex.Lock()
	defer mutex.Unlock()
	c.Close()
}



func SetupConsumer() *Consumer {

	c, err := NewConsumer(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
		"GROUP_ID":       "TEST_GROUP",
	})
	if err != nil {
		panic(err)
	}
	return c
}

func TeardownConsumer(c *Consumer) {
	mutex.Lock()
	defer mutex.Unlock()
	c.Close()
}



func SetupProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
	})
	if err != nil {
		panic(err)
	}
	return p
}

func TeardownProducer(p *kafka.Producer) {
	mutex.Lock()
	defer mutex.Unlock()
	p.Close()
}



func TeardownEnvironment() {
	c := SetupConfigurator()
	defer TeardownConfigurator(c)

	if err := c.DeleteAllTopics(context.Background()); err != nil {
		panic(err)
	}

}



func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}