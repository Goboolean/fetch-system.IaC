package preparer_test

import (
	"context"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.IaC/cmd/wire"
	"github.com/Goboolean/fetch-system.IaC/internal/connect"
	"github.com/Goboolean/fetch-system.IaC/internal/etcd"
	"github.com/Goboolean/fetch-system.IaC/internal/kafka"
	"github.com/Goboolean/fetch-system.IaC/internal/preparer"
	"github.com/stretchr/testify/assert"
)



func TestScenario(t *testing.T) {

	var preparer *preparer.Manager
	var etcd *etcd.Client
	var connect *connect.Client
	var conf *kafka.Configurator
	var topics []string

	const (
		connectorName = "preparer.connector"
		connectorTasks = 10
	)

	t.Run("Setup prerarer", func(tt *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		var err error
		var cleanup func()

		preparer, cleanup, err = wire.InitializePreparer(ctx)
		assert.NoError(tt, err)
		t.Cleanup(cleanup)

		etcd, cleanup, err = wire.InitializeETCDClient(ctx)
		assert.NoError(tt, err)
		t.Cleanup(cleanup)

		connect, cleanup, err = wire.InitializeKafkaConnectClient(ctx)
		assert.NoError(tt, err)
		t.Cleanup(cleanup)

		conf, cleanup, err = wire.InitializeKafkaConfigurator(ctx)
		assert.NoError(tt, err)
		t.Cleanup(cleanup)
	})

	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
		defer cancel()

		err := connect.DeleteAllConnectors(ctx)
		assert.NoError(t, err)

		connectors, err := connect.GetConnectors(ctx)
		assert.NoError(t, err)
		assert.Empty(t, connectors)

		err = etcd.DeleteAllProducts(ctx)
		assert.NoError(t, err)

		err = conf.DeleteAllTopics(ctx)
		assert.NoError(t, err)
	})

	t.Run("SyncETCDToDB", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
		defer cancel()

		var err error
		topics, err = preparer.SyncETCDToDB(ctx)
		assert.NoError(t, err)
		assert.NotEmpty(t, topics)

		products, err := etcd.GetAllProducts(ctx)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(products), 2000)
	})

	t.Run("Prepare", func(t *testing.T) {
		topics = topics[:300]

		start := time.Now()

		ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)
		defer cancel()

		err := preparer.PrepareTopicsBatch(ctx, connectorName, connectorTasks, topics)
		assert.NoError(t, err)

		elasped := time.Since(start)
		t.Log("Elapsed time:", elasped)
	})
}