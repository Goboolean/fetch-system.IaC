package prepare

import (
	"context"
	"fmt"

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




func (m *Manager) SyncETCDToDB(ctx context.Context) ([]string, error) {
	
	products, err := m.db.GetAllProducts(ctx)
	if err != nil {
		return nil, err
	}

	dtos := make([]*etcd.Product, len(products))

	for i, product := range products {
		dtos[i] = &etcd.Product{
			ID:     product.ID,
			Symbol: product.Symbol,
		}
	}

	if err := m.etcd.InsertProducts(ctx, dtos); err != nil {
		return nil, err
	}

	var topics []string
	for _, product := range products {
		topics = append(topics, product.ID)
	}

	return topics, nil
}



func (m *Manager) PrepareTopic(ctx context.Context, topic string) error {

	topicTic := fmt.Sprintf("%s.t", topic)
	topic1s := fmt.Sprintf("%s.1s", topic)
	topic5s := fmt.Sprintf("%s.5s", topic)
	topic1m := fmt.Sprintf("%s.1m", topic)
	topic5m := fmt.Sprintf("%s.5m", topic)

	topicAll := []string{topicTic, topic1s, topic5s, topic1m, topic5m}
	topicAggs := []string{topic1s, topic5s, topic1m, topic5m}

	if err := m.conf.CreateTopics(ctx, topicAll...); err != nil {
		return err
	}

	for _, topic := range topicAggs {
		if err := m.connect.CreateConnector(ctx, topic); err != nil {
			return err
		}
	}

	return nil
}