package influx_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.IaC/pkg/influx"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stretchr/testify/assert"
)

var testClient *influx.DB
var rawClient influxdb2.Client
var testStockID = "stock.aapl.usa"
var testTimeFrame = "1m"

// var options = influx.Opts{
// 	Url:             os.Getenv("INFLUXDB_URL"),
// 	Token:           os.Getenv("INFLUXDB_TOKEN"),
// 	TradeBucketName: os.Getenv("INFLUXDB_TRADE_BUCKET"),
// 	Org:             os.Getenv("INFLUXDB_ORG"),
// }

var options = influx.Opts{
	Url:             "http://localhost:8086",
	Token:           "tokenforadmin",
	TradeBucketName: "bucket",
	Org:             "org",
}

func TestMain(m *testing.M) {
	var err error
	testClient, err = influx.NewDB(&options)
	if err != nil {
		panic(err)
	}
	rawClient = influxdb2.NewClient(options.Url, options.Token)
	err = createBucketIfNotExits(rawClient, options.Org, options.TradeBucketName)
	if err != nil {
		panic(err)
	}

	m.Run()
}

func createBucketIfNotExits(client influxdb2.Client, orgName, bucketName string) error {
	org, err := client.OrganizationsAPI().FindOrganizationByName(context.Background(), orgName)
	if err != nil {
		return err
	}

	_, err = client.BucketsAPI().FindBucketByName(context.Background(), bucketName)

	if err != nil {
		_, err = client.BucketsAPI().CreateBucketWithName(context.Background(), org, bucketName)
		return err
	}
	return nil
}

func RecreateBucket(client influxdb2.Client, orgName, bucketName string) error {

	org, err := client.OrganizationsAPI().FindOrganizationByName(context.Background(), orgName)
	if err != nil {
		return err
	}

	bucket, err := client.BucketsAPI().FindBucketByName(context.Background(), bucketName)
	if err != nil {
		return nil
	}

	client.BucketsAPI().DeleteBucket(context.Background(), bucket)
	_, err = client.BucketsAPI().CreateBucketWithName(context.Background(), org, bucketName)

	return err
}

func TestConstructor(t *testing.T) {
	t.Run("constructor로 생성한 client의 Ping()을 호출했을 때 에러가 없어야 한다.", func(t *testing.T) {
		//act
		err := testClient.Ping(context.Background())
		assert.NoError(t, err)
	})
}

func TestFetchByTimeRange(t *testing.T) {
	t.Run("빈 bucket에 쿼리를 요청한 경우,"+
		"data가 비어있어야 하고 에러가 없어야 한다.", func(t *testing.T) {
		//arrange
		RecreateBucket(rawClient, options.Org, options.TradeBucketName)
		//act
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
		now := time.Now()

		defer cancel()
		aggregates, err := testClient.FetchByTimeRange(ctx, testStockID, "1m", now.Add(-5*60*time.Second), now)
		assert.NoError(t, err)
		assert.Len(t, aggregates, 0)
	})

	t.Run("from과 to가 저장된 데이터의 범위 바깥을 나타내는 경우"+
		"data가 비어있어야 하고 에러가 없어야 한다", func(t *testing.T) {
		//arrange

		RecreateBucket(rawClient, options.Org, options.TradeBucketName)
		writer := rawClient.WriteAPI(options.Org, options.TradeBucketName)
		start := time.Now()
		for i := 0; i < 60; i++ {
			writer.WritePoint(
				write.NewPoint(
					fmt.Sprintf("%s.%s", testStockID, "1m"),
					map[string]string{},
					map[string]interface{}{
						"start":  float64(1.0),
						"stop":   float64(2.0),
						"high":   float64(3.0),
						"low":    float64(4.0),
						"volume": int64(4),
					},
					start.Add(time.Duration(i)*time.Second),
				),
			)
		}
		//act
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
		defer cancel()

		aggregates, err := testClient.FetchByTimeRange(ctx,
			testStockID,
			"1m",
			start.Add(-5*time.Second),
			start.Add(-1*time.Second))

		assert.NoError(t, err)
		assert.Len(t, aggregates, 0)
	})

	t.Run("from과 to가 저장된 데이터를 포함하는 경우,"+
		"알맞은 개수의 데이터가 있어야 하고 에러가 없어야 한다", func(t *testing.T) {
		//arrange

		RecreateBucket(rawClient, options.Org, options.TradeBucketName)
		writer := rawClient.WriteAPIBlocking(options.Org, options.TradeBucketName)
		start := time.Now()
		for i := 0; i < 60; i++ {
			writer.WritePoint(
				context.Background(),
				write.NewPoint(
					fmt.Sprintf("%s.%s", testStockID, "1m"),
					map[string]string{},
					map[string]interface{}{
						"open":   float64(1.0),
						"close":  float64(2.0),
						"high":   float64(3.0),
						"low":    float64(4.0),
						"volume": int64(4),
					},
					start.Add(time.Duration(i)*time.Second),
				),
			)
		}

		//act
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Second))
		defer cancel()

		aggregates, err := testClient.FetchByTimeRange(ctx,
			testStockID,
			"1m",
			start,
			start.Add(30*time.Second),
		)

		assert.NoError(t, err)
		assert.Len(t, aggregates, 30)
	})
}
