package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Goboolean/fetch-system.IaC/cmd/wire"
)




func main() {

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	preparer, cleanup, err := wire.InitializePreparer()
	if err != nil {
		panic(err)
	}
	defer cleanup()

	topics, err := preparer.SyncETCDToDB(ctx)
	if err != nil {
		panic(err)
	}

	if err := preparer.PrepareTopics(ctx, topics); err != nil {
		panic(err)
	}
}