package kafka

import (
	"os"
	"sync"
	"testing"

	_ "github.com/Goboolean/common/pkg/env"
	"github.com/Goboolean/common/pkg/resolver"
	_ "github.com/Goboolean/common/pkg/env"
	"github.com/Goboolean/fetch-system.IaC/internal/kafka"
	log "github.com/sirupsen/logrus"

)

var mutex = &sync.Mutex{}

var conf *kafka.Configurator



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



func TestMain(m *testing.M) {
	conf = SetupConfigurator()
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetLevel(log.TraceLevel)

	code := m.Run()
	os.Exit(code)
	TeardownConfigurator(conf)
}