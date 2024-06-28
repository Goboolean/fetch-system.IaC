package influx

import (
	"context"
	"fmt"

	"github.com/Goboolean/common/pkg/resolver"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type DB struct {
	client           influxdb2.Client
	tradeBucket      string
	reader           api.QueryAPI
	orderWriter      api.WriteAPIBlocking
	annotationWriter api.WriteAPIBlocking
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
	if BucketExists(client, tradeBucket) {
		return nil, fmt.Errorf("influx client: bucket %s does not exist", tradeBucket)
	}

	orderEventBucket, err := c.GetStringKey("INFLUX_ORDER_EVENT_BUCKET")
	if err != nil {
		return nil, fmt.Errorf("influx client: fail to create client %v", err)
	}
	if BucketExists(client, tradeBucket) {
		return nil, fmt.Errorf("influx client: bucket %s does not exist", orderEventBucket)
	}

	annotationBucket, err := c.GetStringKey("INFLUX_ANNOTATION_EVENT_BUCKET")
	if err != nil {
		return nil, fmt.Errorf("influx client: fail to create client %v", err)
	}
	if BucketExists(client, tradeBucket) {
		return nil, fmt.Errorf("influx client: bucket %s does not exist", annotationBucket)
	}

	instance := &DB{
		client:           client,
		tradeBucket:      tradeBucket,
		orderWriter:      client.WriteAPIBlocking(org, orderEventBucket),
		annotationWriter: client.WriteAPIBlocking(org, annotationBucket),
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

func BucketExists(c influxdb2.Client, bucket string) bool {
	bucketApi := c.BucketsAPI()
	_, err := bucketApi.FindBucketByName(context.Background(), bucket)
	return err == nil
}
