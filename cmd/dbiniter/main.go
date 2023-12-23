package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Goboolean/fetch-system.IaC/cmd/wire"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	_ "github.com/Goboolean/common/pkg/env"
	_ "github.com/Goboolean/fetch-system.IaC/internal/log"
)


func main() {
	log.Info("Application started")

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	retriever, cleanup, err := wire.InitializeRetriever(ctx)
	if err != nil {
		log.Panic(errors.Wrap(err, "Failed to initialize retriever"))
	}
	defer cleanup()

	stored, err := retriever.CheckKORStockStored(ctx)
	if err != nil {
		log.Panic(errors.Wrap(err, "Failed to check KOR stock stored"))
	}

	if !stored {
		if err := retriever.StoreKORStocks(ctx); err != nil {
			log.Panic(errors.Wrap(err, "Failed to store KOR stocks"))
		}
	}

	stored, err = retriever.CheckUSAStockStored(ctx)
	if err != nil {
		log.Panic(errors.Wrap(err, "Failed to check USA stock stored"))
	}

	if !stored {
		if err := retriever.StoreUSAStocks(ctx); err != nil {
			log.Panic(errors.Wrap(err, "Failed to store USA stocks"))
		}
	}

	log.Info("Application successfully finished")
}