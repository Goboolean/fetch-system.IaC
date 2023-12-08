package kafka_test

import (
	"os"
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"

	_ "github.com/Goboolean/common/pkg/env"
)

var mutex = &sync.Mutex{}



func TestMain(m *testing.M) {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetLevel(log.TraceLevel)

	code := m.Run()
	os.Exit(code)
}