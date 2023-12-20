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


const batchSize = 500

func (m *Manager) PrepareTopics(ctx context.Context, topics []string) error {
	
	for i := 0; i < len(topics); i += batchSize {
		end := i + batchSize
		if end > len(topics) {
			end = len(topics)
		}

		var connectorName = fmt.Sprintf("%s.%d", "preparer.connector", i)
		var connectorTasks = 10

		if err := m.PrepareTopicsBatch(ctx, connectorName, connectorTasks, topics[i:end]); err != nil {
			return err
		}
	}

	return nil
}


func (m *Manager) PrepareTopicsBatch(ctx context.Context, connectorName string, connectorTasks int, topics []string) error {

	topicAll := make([]string, 0)
	topicAggs := make([]string, 0)

	for _, topic := range topics {
		topicTic := fmt.Sprintf("%s.t", topic)
		topic1s := fmt.Sprintf("%s.1s", topic)
		topic5s := fmt.Sprintf("%s.5s", topic)
		topic1m := fmt.Sprintf("%s.1m", topic)
		topic5m := fmt.Sprintf("%s.5m", topic)

		topicAll = append(topicAll, topicTic, topic1s, topic5s, topic1m, topic5m)
		topicAggs = append(topicAggs, topic1s, topic5s, topic1m, topic5m)
	}

	if err := m.conf.CreateTopics(ctx, topicAll...); err != nil {
		return err
	}

	configs := make([]connect.ConnectorTopicConfig, len(topicAggs))
	for i, topic := range topicAggs {
		configs[i] = connect.ConnectorTopicConfig{
			Topic:            topic,
			Collection:       topic,
			RotateIntervalMs: 100000,
		}
	}

	if err := m.connect.CreateBulkTopicConnector(ctx, connectorName, connectorTasks, configs); err != nil {
		return err
	}

	_, err := m.connect.CheckTasksStatus(ctx, connectorName)
	if err != nil {
		return err
	}
	return nil
}