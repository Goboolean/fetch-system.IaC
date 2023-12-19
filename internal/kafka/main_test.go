package kafka_test

import (
	"os"
	"sync"

	_ "github.com/Goboolean/common/pkg/env"
	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/kafka"
)

var mutex = &sync.Mutex{}




func SetupConfigurator() *kafka.Configurator {

	conf, err := kafka.New(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
	})
	if err != nil {
		panic(err)
	}

	return conf
}

func TeardownConfigurator(conf *kafka.Configurator) {
	mutex.Lock()
	defer mutex.Unlock()
	conf.Close()
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