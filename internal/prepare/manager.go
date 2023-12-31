package prepare

import (
	"context"
	"fmt"
	"time"

	"github.com/Goboolean/fetch-system.IaC/internal/connect"
	"github.com/Goboolean/fetch-system.IaC/internal/etcd"
	"github.com/Goboolean/fetch-system.IaC/internal/kafka"
	"github.com/Goboolean/fetch-system.IaC/pkg/db"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
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
	log.Infof("received number of %d products", len(products))

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

func (m *Manager) PrepareTopics(ctx context.Context, baseConnectorName string, topics []string) error {
	log.Info("Preparing topics started")

	for i := 0; i < len(topics); i += batchSize {
		end := i + batchSize
		if end > len(topics) {
			end = len(topics)
		}

		var connectorName = fmt.Sprintf("%s.[%d:%d]", baseConnectorName, i, end)
		var connectorTasks = 10

		start := time.Now()
		log.Infof("Preparing topics batch started: %s)", connectorName)

		if err := m.PrepareTopicsBatch(ctx, connectorName, connectorTasks, topics[i:end]); err != nil {
			return err
		}

		log.Infof("Preparing topics batch took: %s)", time.Since(start))
	}

	log.Info("Preparing topics finished")
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
		if errors.Is(err, kafka.ErrSomeOfTopicAlreadyExist) {
			log.Warn("Topic already exists")
		} else {
			return errors.Wrap(err, "Failed to create topics")
		}
	} else {
		log.Info("Topics are successfully created")
	}

	configs := make([]connect.ConnectorTopicConfig, len(topicAggs))
	for i, topic := range topicAggs {
		configs[i] = connect.ConnectorTopicConfig{
			Topic:            topic,
			Collection:       topic,
			RotateIntervalMs: 100000,
		}
	}

	exists, err := m.connect.CheckConnectorExists(ctx, connectorName)
	if err != nil {
		return errors.Wrap(err, "Failed to check tasks status")
	}

	if exists {
		log.Warn("Connector already exists")
		return nil
	}

	if err := m.connect.CreateBulkTopicConnector(ctx, connectorName, connectorTasks, configs); err != nil {
		return errors.Wrap(err, "Failed to create connector")
	}

	_, err = m.connect.CheckTasksStatus(ctx, connectorName)
	if err != nil {
		return errors.Wrap(err, "Failed to check tasks status")
	}
	return nil
}