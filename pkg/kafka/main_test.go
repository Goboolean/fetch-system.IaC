package kafka_test

import (
	"os"
	"sync"
	"testing"

	_ "github.com/Goboolean/common/pkg/env"
	log "github.com/sirupsen/logrus"
)

var mutex = &sync.Mutex{}

// List for perfect test coverage
// [*] 1. Ping test
// [*] 2. Produce and consume
// [ ] 3. Produce and consume not existing topic
// [*] 4. Consume with different group
// [*] 5. Consume with same group
// [ ] 6. Consume with registry
// [ ] 7. Consume invalid typed message



func TestMain(m *testing.M) {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})
	log.SetLevel(log.TraceLevel)

	code := m.Run()
	os.Exit(code)
}