package influx

import (
	"context"
	"fmt"

	"github.com/Goboolean/common/pkg/resolver"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	log "github.com/sirupsen/logrus"
)

const tradeBucket = "minimal"
const orderEventBucket = "order"
const annotationBucket = "annotation"

type DB struct {
	client           influxdb2.Client
	reader           api.QueryAPI
	orderWriter      api.WriteAPI
	annotationWriter api.WriteAPI
}

func NewDB(c *resolver.ConfigMap) (*DB, error) {

	url, err := c.GetStringKey("INFLUX_URL")
	if err != nil {
		return nil, fmt.Errorf("influx client: fail to create client %V", err)
	}

	token, err := c.GetStringKey("INFLUX_TOKEN")
	if err != nil {
		return nil, fmt.Errorf("influx client: fail to create client %v", err)
	}

	org, err := c.GetStringKey("INFLUX_ORG")
	if err != nil {
		return nil, fmt.Errorf("influx client: fail to create client %v", err)
	}

	client := influxdb2.NewClient(url, token)
	instance := &DB{
		client:           client,
		orderWriter:      client.WriteAPI(org, orderEventBucket),
		annotationWriter: client.WriteAPI(org, annotationBucket),
	}

	go func() {
		for e := range instance.orderWriter.Errors() {
			log.Error(e)
		}
	}()

	go func() {
		for e := range instance.annotationWriter.Errors() {
			log.Error(e)
		}
	}()

	return instance, nil
}

func (d *DB) Flush(ctx context.Context) {
	d.orderWriter.Flush()
	d.annotationWriter.Flush()
}

func (d *DB) Ping(ctx context.Context) error {
	_, err := d.client.Ping(ctx)
	return err
}

func (d *DB) Close() error {
	d.client.Close()
	return nil
}
