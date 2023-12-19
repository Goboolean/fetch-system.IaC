package prepare_test

import (
	"context"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.IaC/cmd/wire"
	"github.com/Goboolean/fetch-system.IaC/internal/etcd"
	"github.com/Goboolean/fetch-system.IaC/internal/prepare"
	"github.com/stretchr/testify/assert"
)



func TestScenario(t *testing.T) {

	var preparer *prepare.Manager
	var etcd *etcd.Client
	var topics []string

	t.Run("Setup prerarer", func(tt *testing.T) {
		var err error
		var cleanup func()

		preparer, cleanup, err = wire.InitializePreparer()
		assert.NoError(tt, err)
		t.Cleanup(cleanup)

		etcd, cleanup, err = wire.InitializeETCDClient()
		assert.NoError(tt, err)
		t.Cleanup(cleanup)
	})

	t.Run("SyncETCDToDB", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
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
		start := time.Now()

		for _, topic := range topics {
			ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
			defer cancel()

			err := preparer.PrepareTopic(ctx, topic)
			assert.NoError(t, err)
		}

		elasped := time.Since(start)
		t.Log("Elapsed time:", elasped)
	})

	t.Run("CheckPrepared", func(t *testing.T) {
		connect, cleanup, err := wire.InitializeKafkaConnectClient()
		assert.NoError(t, err)
		t.Cleanup(cleanup)

		ctx, cancel := context.WithTimeout(context.Background(), 3 * time.Second)
		defer cancel()

		for _, topic := range topics {
			err := connect.CheckPluginConfig(ctx, topic)
			assert.NoError(t, err)
		}
	})
}