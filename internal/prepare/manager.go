package prepare

import (
	"context"

	"github.com/Goboolean/fetch-system.IaC/internal/connect"
	"github.com/Goboolean/fetch-system.IaC/internal/etcd"
	"github.com/Goboolean/fetch-system.IaC/internal/kafka"
	"github.com/Goboolean/fetch-system.IaC/pkg/db"
)




type Manager struct {
	etcd    *etcd.Client
	db      *db.Client
	connect *connect.Client
	conf    *kafka.Configurator
}


func New(
	etcd    *etcd.Client,
	db      *db.Client,
	connect *connect.Client,
	conf    *kafka.Configurator) *Manager {
		return &Manager{
			etcd:    etcd,
			db:      db,
			connect: connect,
			conf:    conf,
		}
}




func (m *Manager) SyncETCDToDB(ctx context.Context) error {
	
	products, err := m.db.GetAllProducts(ctx)
	if err != nil {
		return err
	}

	dtos := make([]*etcd.Product, len(products))

	for i, product := range products {
		dtos[i] = &etcd.Product{
			ID:     product.ID,
			Symbol: product.Symbol,
			Type:   product.Market,
		}
	}

	if err := m.etcd.InsertProducts(ctx, dtos); err != nil {
		return err
	}

	return nil
}