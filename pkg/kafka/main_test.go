package kafka_test

import (
	"os"
	"sync"
	"testing"

	_ "github.com/Goboolean/common/pkg/env"
	"github.com/Goboolean/common/pkg/resolver"
	_ "github.com/Goboolean/common/pkg/env"
	"github.com/Goboolean/fetch-system.IaC/internal/kafkaadmin"
	log "github.com/sirupsen/logrus"

)

var mutex = &sync.Mutex{}

var conf *kafkaadmin.Configurator



func SetupConfigurator() *kafkaadmin.Configurator {
	c, err := kafkaadmin.New(&resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),
	})
	if err != nil {
		panic(err)
	}
	return c
}

func TeardownConfigurator(c *kafkaadmin.Configurator) {
	mutex.Lock()
	defer mutex.Unlock()
	c.Close()
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