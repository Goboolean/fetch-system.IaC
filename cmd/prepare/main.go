package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/Goboolean/common/pkg/env"
	"github.com/Goboolean/fetch-system.IaC/cmd/wire"
	_ "github.com/Goboolean/fetch-system.IaC/internal/log"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)




func main() {
	log.Info("Application started")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	preparer, cleanup, err := wire.InitializePreparer(ctx)
	if err != nil {
		log.Panic(errors.Wrap(err, "Failed to initialize preparer"))
	}
	defer cleanup()

	topics, err := preparer.SyncETCDToDB(ctx)
	if err != nil {
		log.Panic(errors.Wrap(err, "Failed to synchronize etcd to db"))
	}

	if err := preparer.PrepareTopics(ctx, topics); err != nil {
		log.Panic(errors.Wrap(err, "Failed to prepare topics"))
	}

	log.Info("Application successfully finished")
}