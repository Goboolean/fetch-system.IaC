package connect_test

import (
	"os"
	"sync"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/pkg/kafka"
	"github.com/Goboolean/fetch-system.IaC/pkg/mongo"
	"github.com/Goboolean/fetch-system.IaC/internal/connect"

	_ "github.com/Goboolean/common/pkg/env"
)

var mutex = &sync.Mutex{}



func SetupConnect() *connect.Client {
	c, err := connect.New(&resolver.ConfigMap{
		"HOST": os.Getenv("KAFKA_CONNECT_HOST"),
		"MONGODB_CONNECTION_URI": os.Getenv("MONGODB_CONNECTION_URI"),
		"MONGODB_DATABASE": os.Getenv("MONGODB_DATABASE"),
	})
	if err != nil {
		panic(err)
	}
	return c
}

func TeardownConnect(c *connect.Client) {
	c.Close()
}



func SetupProducer() *kafka.Producer {
	p, err := kafka.NewProducer(&resolver.ConfigMap{
		"KAFKA_BROKER": os.Getenv("KAFKA_BROKER"),
		"KAFKA_TOPIC": os.Getenv("KAFKA_TOPIC"),
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


func SetupConsumer() *kafka.Consumer {
	c, err := kafka.NewConsumer(&resolver.ConfigMap{
		"KAFKA_BROKER": os.Getenv("KAFKA_BROKER"),
		"KAFKA_TOPIC": os.Getenv("KAFKA_TOPIC"),
	})
	if err != nil {
		panic(err)
	}
	return c
}

func TeardownConsumer(c *kafka.Consumer) {
	mutex.Lock()
	defer mutex.Unlock()
	c.Close()
}



func SetupMongoClient() *mongo.DB {
	c, err := mongo.NewDB(&resolver.ConfigMap{
		"MONGODB_CONNECTION_URI": os.Getenv("MONGODB_CONNECTION_URI"),
		"MONGODB_DATABASE": os.Getenv("MONGODB_DATABASE"),
	})
	if err != nil {
		panic(err)
	}
	return c
}

func TeardownMongoClient(c *mongo.DB) {
	c.Close()
}


