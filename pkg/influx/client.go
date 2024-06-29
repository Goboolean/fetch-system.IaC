package influx

import (
	"context"
	"fmt"

	"github.com/Goboolean/common/pkg/resolver"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type DB struct {
	client      influxdb2.Client
	tradeBucket string
	reader      api.QueryAPI
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

	tradeBucket, err := c.GetStringKey("INFLUX_TRADE_BUCKET")
	if err != nil {
		return nil, fmt.Errorf("influx client: fail to create client %v", err)
	}
	if bucketExists(client, tradeBucket) {
		return nil, fmt.Errorf("influx client: bucket %s does not exist", tradeBucket)
	}

	instance := &DB{
		client:      client,
		tradeBucket: tradeBucket,
		reader:      client.QueryAPI(org),
	}

	return instance, nil
}

func (d *DB) Ping(ctx context.Context) error {
	_, err := d.client.Ping(ctx)
	return err
}

func (d *DB) Close() error {
	d.client.Close()
	return nil
}

func bucketExists(c influxdb2.Client, bucket string) bool {
	bucketApi := c.BucketsAPI()
	_, err := bucketApi.FindBucketByName(context.Background(), bucket)
	return err == nil
}
