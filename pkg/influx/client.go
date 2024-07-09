package influx

import (
	"context"
	"fmt"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

type DB struct {
	client      influxdb2.Client
	tradeBucket string
	reader      api.QueryAPI
}

type Opts struct {
	URL             string
	Token           string
	Org             string
	TradeBucketName string
}

func NewDB(o *Opts) (*DB, error) {

	if o.URL == "" {
		return nil, fmt.Errorf("create influx db client: Required field Url is blank")
	}

	if o.Token == "" {
		return nil, fmt.Errorf("create influx db client: Required field Token is blank")
	}

	if o.Org == "" {
		return nil, fmt.Errorf("create influx db client: Required field Url is blank")
	}

	if o.TradeBucketName == "" {
		return nil, fmt.Errorf("create influx db client: Required field TradeBucketName is blank")
	}

	client := influxdb2.NewClient(o.URL, o.Token)

	instance := &DB{
		client:      client,
		tradeBucket: o.TradeBucketName,
		reader:      client.QueryAPI(o.Org),
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
