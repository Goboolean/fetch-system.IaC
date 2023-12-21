package mongo

import (
	"context"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client    *mongo.Client
	defaultDB string
}

func NewDB(c *resolver.ConfigMap) (*DB, error) {

	conn_url, err := c.GetStringKey("CONNECTION_URI")
	if err != nil {
		return nil, err
	}

	database, err := c.GetStringKey("DATABASE")
	if err != nil {
		return nil, err
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(conn_url).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)

	if err != nil {
		return nil, err
	}

	return &DB{
		client:    client,
		defaultDB: database,
	}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := db.client.Ping(ctx, nil); err != nil {
			time.Sleep(time.Second)
			continue
		}

		return nil
	}
}

func (db *DB) Close() error {
	return db.client.Disconnect(context.Background())
}