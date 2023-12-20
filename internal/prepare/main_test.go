package prepare_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.IaC/cmd/wire"
)


var cleanups []func()

func Setup() {
	retriever, cleanup, err := wire.InitializeRetriever()
	if err != nil {
		panic(err)
	}
	cleanups = append(cleanups, cleanup)

	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	if err := retriever.StoreKORStocks(ctx); err != nil {
		panic(err)
	}
}


func Teardown() {
	db, cleanup, err := wire.InitializePostgreSQLClient()
	if err != nil {
		panic(err)
	}
	cleanups = append(cleanups, cleanup)

	ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
	defer cancel()

	if err := db.DeleteAllProducts(ctx); err != nil {
		panic(err)
	}

	etcd, cleanup, err := wire.InitializeETCDClient()
	if err != nil {
		panic(err)
	}
	cleanups = append(cleanups, cleanup)

	if etcd.DeleteAllProducts(ctx); err != nil {
		panic(err)
	}

	for _, cleanup := range cleanups {
		cleanup()
	}
}


func TestMain(m *testing.M) {
	Setup()
	code := m.Run()
	Teardown()
	os.Exit(code)
}