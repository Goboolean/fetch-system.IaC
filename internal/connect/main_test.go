package connect_test

import (
	"os"
	"sync"
	"testing"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/connect"
	"github.com/Goboolean/fetch-system.IaC/internal/kafka"
	_kafka "github.com/Goboolean/fetch-system.IaC/pkg/kafka"
	"github.com/Goboolean/fetch-system.IaC/pkg/mongo"

	_ "github.com/Goboolean/common/pkg/env"
)

var mutex = sync.Mutex{}
var conf *kafka.Configurator


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



func TeardownConsumer(c *_kafka.Consumer) {
	mutex.Lock()
	defer mutex.Unlock()
	c.Close()
}


func SetupAdminClient() *kafka.Configurator {
	a, err := kafka.New(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
	})
	if err != nil {
		panic(err)
	}
	return a
}

func TeardownAdminClient(a *kafka.Configurator) {
	mutex.Lock()
	a.Close()
	mutex.Unlock()
}


func SetupMongoClient() *mongo.DB {
	c, err := mongo.NewDB(&resolver.ConfigMap{
		"CONNECTION_URI": os.Getenv("MONGODB_CONNECTION_URI"),
		"DATABASE": os.Getenv("MONGODB_DATABASE"),
	})
	if err != nil {
		panic(err)
	}
	return c
}

func TeardownMongoClient(c *mongo.DB) {
	c.Close()
}



func TestMain(m *testing.M) {

	conf = SetupAdminClient()
	code := m.Run()
	os.Exit(code)
	TeardownAdminClient(conf)
}