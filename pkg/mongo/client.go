package mongo

import (
	"context"
	"fmt"
	"net/url"

	"github.com/Goboolean/common/pkg/resolver"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client    *mongo.Client
	defaultDB string
}

func NewDB(c *resolver.ConfigMap) (*DB, error) {

	user, err := c.GetStringKey("USER")
	if err != nil {
		return nil, err
	}

	password, err := c.GetStringKey("PASSWORD")
	if err != nil {
		return nil, err
	}

	host, err := c.GetStringKey("HOST")
	if err != nil {
		return nil, err
	}

	port, err := c.GetStringKey("PORT")
	if err != nil {
		return nil, err
	}

	database, err := c.GetStringKey("DATABASE")
	if err != nil {
		return nil, err
	}

	address := fmt.Sprintf("%s:%s", host, port)

	u := &url.URL{
		Scheme: "mongodb",
		User:   url.UserPassword(user, password),
		Host:   address,
		Path:   "/",
		RawQuery: url.Values{
			"authSource": []string{database},
		}.Encode(),
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(u.String()).SetServerAPIOptions(serverAPI)

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
	return db.client.Ping(ctx, nil)
}

func (db *DB) Close() error {
	return db.client.Disconnect(context.Background())
}