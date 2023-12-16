package task

import (
	"context"
	"fmt"

	"github.com/Goboolean/fetch-system.IaC/internal/connect"
	"github.com/Goboolean/fetch-system.IaC/internal/etcd"
	"github.com/Goboolean/fetch-system.IaC/internal/kafka"
	"github.com/Goboolean/fetch-system.IaC/internal/kis"
	"github.com/Goboolean/fetch-system.IaC/internal/polygon"
	"github.com/Goboolean/fetch-system.IaC/internal/rdbms"
)




type Manager struct {
	polygon *polygon.Client
	etcd    *etcd.Client
	db      *rdbms.Client
	kis     *kis.Reader
	connect *connect.Client
	conf    *kafka.Configurator
}


func New(
	polygon *polygon.Client,
	etcd    *etcd.Client,
	db      *rdbms.Client,
	kis     *kis.Reader,
	connect *connect.Client,
	conf    *kafka.Configurator) *Manager {
		return &Manager{
			polygon: polygon,
			etcd:    etcd,
			db:      db,
			kis:     kis,
			connect: connect,
			conf:    conf,
		}
}


func (m *Manager) InitUSAStocks(ctx context.Context) error {

	details, err := m.GetAllUSATickerDetails(ctx)
	if err != nil {
		return err
	}

	dtos := make([]rdbms.InsertProductsParams, len(details))

	for i, detail := range details {
		dtos[i] = rdbms.InsertProductsParams{
			ID:     fmt.Sprintf("%s.%s.%s", "stock", detail.Ticker, "usa"),
			Symbol: detail.Ticker,
			Locale: "usa",
			Market: "stock",
		}
	}

	count, err := m.db.InsertProducts(ctx, dtos)
	if err != nil {
		return err
	}

	if int(count) != len(details) {
		return fmt.Errorf("failed to insert products")
	}
	return nil
}


func (m *Manager) InitKORStocks(ctx context.Context) error {

	details, err := m.kis.ReadAllTickerDetalis("./api/csv/data.csv")
	if err != nil {
		return err
	}

	dtos := make([]rdbms.InsertProductsParams, len(details))

	for i, detail := range details {
		dtos[i] = rdbms.InsertProductsParams{
			ID:     fmt.Sprintf("%s.%s.%s", "stock", detail.Ticker, "kor"),
			Symbol: detail.Ticker,
			Locale: "kor",
			Market: "stock",
		}
	}

	count, err := m.db.InsertProducts(ctx, dtos)
	if err != nil {
		return err
	}

	if int(count) != len(details) {
		return fmt.Errorf("failed to insert products")
	}
	return nil
}


