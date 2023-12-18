//go:build wireinject
// +build wireinject

package wire

import (
	"os"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/connect"
	"github.com/Goboolean/fetch-system.IaC/internal/etcd"
	"github.com/Goboolean/fetch-system.IaC/internal/kafka"
	"github.com/Goboolean/fetch-system.IaC/internal/kis"
	"github.com/Goboolean/fetch-system.IaC/internal/polygon"
	"github.com/Goboolean/fetch-system.IaC/internal/prepare"
	"github.com/Goboolean/fetch-system.IaC/internal/retrieve"
	"github.com/Goboolean/fetch-system.IaC/pkg/db"
	"github.com/google/wire"

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



func ProvideKafkaConfigurator(c *resolver.ConfigMap) (*kafka.Configurator, func(), error) {
	k, err := kafka.New(c)
	if err != nil {
		return nil, nil, err
	}

	return k, func() {
		k.Close()
	}, nil
}

func ProvideKafkaProducer(c *resolver.ConfigMap) (*kafka.Configurator, func(), error) {
	k, err := kafka.New(c)
	if err != nil {
		return nil, nil, err
	}

	return k, func() {
		k.Close()
	}, nil
}


func ProvideETCDClient(c *resolver.ConfigMap) (*etcd.Client, func(), error) {
	e, err := etcd.New(c)
	if err != nil {
		return nil, nil, err
	}

	return e, func() {
		e.Close()
	}, nil
}

func ProvidePostgreSQLClient(c *resolver.ConfigMap) (*db.Client, func(), error) {
	p, err := db.NewDB(c)
	if err != nil {
		return nil, nil, err
	}

	return p, func() {
		p.Close()
	}, nil
}

func ProvideKafkaConnectClient(c *resolver.ConfigMap) (*connect.Client, func(), error) {
	k, err := connect.New(c)
	if err != nil {
		return nil, nil, err
	}

	return k, func() {
		k.Close()
	}, nil
}

func ProvideKISReader(c *resolver.ConfigMap) (*kis.Reader, error) {
	return kis.New(c)
}

func ProvidePolygonClient(c *resolver.ConfigMap) (*polygon.Client, error) {
	return polygon.New(c)
}



func InitializeKafkaConfigurator() (*kafka.Configurator, func(), error) {
	wire.Build(
		ProvideKafkaConfig,
		ProvideKafkaConfigurator,
	)
	return nil, nil, nil
}

func InitializeKafkaProducer() (*kafka.Configurator, func(), error) {
	wire.Build(
		ProvideKafkaConfig,
		ProvideKafkaProducer,
	)
	return nil, nil, nil
}

func InitializeETCDClient() (*etcd.Client, func(), error) {
	wire.Build(
		ProvideETCDConfig,
		ProvideETCDClient,
	)
	return nil, nil, nil
}

func InitializePostgreSQLClient() (*db.Client, func(), error) {
	wire.Build(
		ProvidePostgreSQLConfig,
		ProvidePostgreSQLClient,
	)
	return nil, nil, nil
}

func InitializeKafkaConnectClient() (*connect.Client, func(), error) {
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



func InitializePreparer() (*prepare.Manager, func(), error) {
	wire.Build(
		InitializeKafkaConfigurator,
		InitializeETCDClient,
		InitializePostgreSQLClient,
		InitializeKafkaConnectClient,
		prepare.New,
	)
	return nil, nil, nil
}

func InitializeRetriever() (*retrieve.Manager, func(), error) {
	wire.Build(
		InitializePostgreSQLClient,
		InitializeKISReader,
		InitializePolygonClient,
		retrieve.New,
	)
	return nil, nil, nil
}