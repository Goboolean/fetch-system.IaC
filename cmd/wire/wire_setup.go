//go:build wireinject
// +build wireinject

package wire

import (
	"context"
	"os"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/connect"
	"github.com/Goboolean/fetch-system.IaC/internal/etcd"
	"github.com/Goboolean/fetch-system.IaC/internal/kafka"
	"github.com/Goboolean/fetch-system.IaC/internal/kis"
	"github.com/Goboolean/fetch-system.IaC/internal/polygon"
	"github.com/Goboolean/fetch-system.IaC/internal/preparer"
	"github.com/Goboolean/fetch-system.IaC/internal/dbiniter"
	"github.com/Goboolean/fetch-system.IaC/pkg/db"
	"github.com/google/wire"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"

	_ "github.com/Goboolean/common/pkg/env"
)



func ProvideKafkaConfig() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"BOOTSTRAP_HOST": os.Getenv("KAFKA_BOOTSTRAP_HOST"),		
	}
}

func ProvideETCDConfig() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"HOST": os.Getenv("ETCD_HOST"),
	}
}

func ProvidePostgreSQLConfig() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"HOST": os.Getenv("POSTGRES_HOST"),
		"PORT": os.Getenv("POSTGRES_PORT"),
		"USER": os.Getenv("POSTGRES_USER"),
		"PASSWORD": os.Getenv("POSTGRES_PASSWORD"),
		"DATABASE": os.Getenv("POSTGRES_DATABASE"),
	}
}

func ProvideKafkaConnectConfig() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"HOST": os.Getenv("KAFKA_CONNECT_HOST"),
		"MONGODB_CONNECTION_URI": os.Getenv("MONGODB_CONNECTION_URI"),
		"MONGODB_DATABASE": os.Getenv("MONGODB_DATABASE"),
	}
}

func ProvideKISConfig() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"FILEPATH": "./api/csv/data.csv",
	}
}

func ProvidePolygonConfig() *resolver.ConfigMap {
	return &resolver.ConfigMap{
		"SECRET_KEY": os.Getenv("POLYGON_SECRET_KEY"),
	}
}



func ProvideKafkaConfigurator(ctx context.Context, c *resolver.ConfigMap) (*kafka.Configurator, func(), error) {
	k, err := kafka.New(c)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to create kafka configurator")
	}
	if err := k.Ping(ctx); err != nil {
		return nil, nil, errors.Wrap(err, "Failed to send ping to kafka configurator")
	}
	log.Info("Kafka configurator is ready")

	return k, func() {
		k.Close()
		log.Info("Kafka configurator is successfully closed")
	}, nil
}

func ProvideKafkaProducer(ctx context.Context, c *resolver.ConfigMap) (*kafka.Configurator, func(), error) {
	k, err := kafka.New(c)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to create kafka producer")
	}
	if err := k.Ping(ctx); err != nil {
		return nil, nil, errors.Wrap(err, "Failed to send ping to kafka producer")
	}
	log.Info("Kafka producer is ready")

	return k, func() {
		k.Close()
		log.Info("Kafka producer is successfully closed")
	}, nil
}


func ProvideETCDClient(ctx context.Context, c *resolver.ConfigMap) (*etcd.Client, func(), error) {
	e, err := etcd.New(c)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to create etcd client")
	}
	if err := e.Ping(ctx); err != nil {
		return nil, nil, errors.Wrap(err, "Failed to send ping to etcd client")
	}
	log.Info("ETCD client is ready")

	return e, func() {
		if err := e.Close(); err != nil {
			log.Error(errors.Wrap(err, "Failed to close etcd client"))
		} else {
			log.Info("ETCD client is successfully closed")
		}
	}, nil
}

func ProvidePostgreSQLClient(ctx context.Context, c *resolver.ConfigMap) (*db.Client, func(), error) {
	p, err := db.NewDB(c)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to create postgresql client")
	}
	if err := p.Ping(ctx); err != nil {
		return nil, nil, errors.Wrap(err, "Failed to send ping to postgresql client")
	}
	log.Info("PostgreSQL client is ready")

	return p, func() {
		p.Close()
		log.Info("PostgreSQL client is successfully closed")
	}, nil
}

func ProvideKafkaConnectClient(ctx context.Context, c *resolver.ConfigMap) (*connect.Client, func(), error) {
	k, err := connect.New(c)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Failed to create kafka connect client")
	}
	if err := k.Ping(ctx); err != nil {
		return nil, nil, errors.Wrap(err, "Failed to send ping to kafka connect client")
	}
	log.Info("Kafka connect client is ready")

	return k, func() {
		k.Close()
		log.Info("Kafka connect client is successfully closed")
	}, nil
}

func ProvideKISReader(c *resolver.ConfigMap) (*kis.Reader, error) {
	return kis.New(c)
}

func ProvidePolygonClient(c *resolver.ConfigMap) (*polygon.Client, error) {
	return polygon.New(c)
}



func InitializeKafkaConfigurator(ctx context.Context) (*kafka.Configurator, func(), error) {
	wire.Build(
		ProvideKafkaConfig,
		ProvideKafkaConfigurator,
	)
	return nil, nil, nil
}

func InitializeKafkaProducer(ctx context.Context) (*kafka.Configurator, func(), error) {
	wire.Build(
		ProvideKafkaConfig,
		ProvideKafkaProducer,
	)
	return nil, nil, nil
}

func InitializeETCDClient(ctx context.Context) (*etcd.Client, func(), error) {
	wire.Build(
		ProvideETCDConfig,
		ProvideETCDClient,
	)
	return nil, nil, nil
}

func InitializePostgreSQLClient(ctx context.Context) (*db.Client, func(), error) {
	wire.Build(
		ProvidePostgreSQLConfig,
		ProvidePostgreSQLClient,
	)
	return nil, nil, nil
}

func InitializeKafkaConnectClient(ctx context.Context) (*connect.Client, func(), error) {
	wire.Build(
		ProvideKafkaConnectConfig,
		ProvideKafkaConnectClient,
	)
	return nil, nil, nil
}

func InitializeKISReader() (*kis.Reader, error) {
	wire.Build(
		ProvideKISConfig,
		ProvideKISReader,
	)
	return nil, nil
}

func InitializePolygonClient() (*polygon.Client, error) {
	wire.Build(
		ProvidePolygonConfig,
		ProvidePolygonClient,
	)
	return nil, nil
}



func InitializePreparer(ctx context.Context) (*preparer.Manager, func(), error) {
	wire.Build(
		InitializeKafkaConfigurator,
		InitializeETCDClient,
		InitializePostgreSQLClient,
		InitializeKafkaConnectClient,
		preparer.New,
	)
	return nil, nil, nil
}

func InitializeRetriever(ctx context.Context) (*dbiniter.Manager, func(), error) {
	wire.Build(
		InitializePostgreSQLClient,
		InitializeKISReader,
		InitializePolygonClient,
		dbiniter.New,
	)
	return nil, nil, nil
}