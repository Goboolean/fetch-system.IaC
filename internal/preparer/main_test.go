package preparer_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.IaC/cmd/wire"
)


var cleanups []func()

func Setup() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbiniter, cleanup, err := wire.InitializeRetriever(ctx)
	if err != nil {
		panic(err)
	}
	cleanups = append(cleanups, cleanup)

	if err := dbiniter.StoreKORStocks(ctx); err != nil {
		panic(err)
	}
}


func Teardown() {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	db, cleanup, err := wire.InitializePostgreSQLClient(ctx)
	if err != nil {
		panic(err)
	}
	cleanups = append(cleanups, cleanup)

	if err := db.DeleteAllProducts(ctx); err != nil {
		panic(err)
	}

	etcd, cleanup, err := wire.InitializeETCDClient(ctx)
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