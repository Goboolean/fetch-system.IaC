package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Goboolean/fetch-system.IaC/cmd/wire"

	_ "github.com/Goboolean/common/pkg/env"
)


func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	manager, cleanup, err := wire.InitializeRetriever()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	if err := manager.StoreKORStocks(ctx); err != nil {
		panic(err)
	}

	if err := manager.StoreUSAStocks(ctx); err != nil {
		panic(err)
	}
}