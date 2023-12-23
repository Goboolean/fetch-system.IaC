package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/Goboolean/common/pkg/env"
	"github.com/Goboolean/fetch-system.IaC/cmd/wire"
	_ "github.com/Goboolean/fetch-system.IaC/internal/log"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)



func main() {
	log.Info("Application started")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 30)

	preparer, cleanup, err := wire.InitializePreparer(ctx)
	if err != nil {
		log.Panic(errors.Wrap(err, "Failed to initialize preparer"))
	}
	defer cleanup()
	cancel()

	ctx, cancel = signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	topics, err := preparer.SyncETCDToDB(ctx)
	if err != nil {
		log.Panic(errors.Wrap(err, "Failed to synchronize etcd to db"))
	}

	if err := preparer.PrepareTopics(ctx, "connector", topics); err != nil {
		log.Panic(errors.Wrap(err, "Failed to prepare topics"))
	}

	log.Info("Application successfully finished")
}