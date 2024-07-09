package influx_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.IaC/pkg/influx"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var testStockID = "stock.aapl.usa"
var testTimeFrame = "1m"

type InfluxTestSuite struct {
	suite.Suite
	testClient *influx.DB
	rawClient  influxdb2.Client
	options    influx.Opts
}

func (suite *InfluxTestSuite) SetupSuite() {
	suite.options = influx.Opts{
		Url:             os.Getenv("INFLUXDB_URL"),
		Token:           os.Getenv("INFLUXDB_TOKEN"),
		TradeBucketName: os.Getenv("INFLUXDB_TRADE_BUCKET"),
		Org:             os.Getenv("INFLUXDB_ORG"),
	}

	var err error
	suite.testClient, err = influx.NewDB(&suite.options)
	if err != nil {
		panic(err)
	}
	suite.rawClient = influxdb2.NewClient(suite.options.Url, suite.options.Token)
	err = suite.createBucketIfNotExits(suite.options.Org, suite.options.TradeBucketName)
	if err != nil {
		panic(err)
	}
}

func (suite *InfluxTestSuite) TestConstructor() {
	// act
	err := suite.testClient.Ping(context.Background())
	assert.NoError(suite.T(), err)

}

func (suite *InfluxTestSuite) SetupSubTest() {
	err := suite.recreateBucket(suite.options.Org, suite.options.TradeBucketName)
	if err != nil {
		suite.FailNow(err.Error())
	}
}

func (suite *InfluxTestSuite) TestFetchByTimeRange() {
	suite.Run("빈 bucket에 쿼리를 요청한 경우, data가 비어있어야 하고 에러가 없어야 한다.", func() {
		// arrange
		// act
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
		now := time.Now()
		defer cancel()

		aggregates, err := suite.testClient.FetchByTimeRange(ctx, testStockID, "1m", now.Add(-5*60*time.Second), now)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), aggregates, 0)
	})

	suite.Run("from과 to가 저장된 데이터의 범위 바깥을 나타내는 경우, data가 비어있어야 하고 에러가 없어야 한다", func() {
		// arrange

		start := time.Now()
		suite.storeStockAggregate(start, time.Second, 60)
		// act
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
		defer cancel()

		aggregates, err := suite.testClient.FetchByTimeRange(ctx, testStockID, testTimeFrame, start.Add(-5*time.Second), start.Add(-1*time.Second))
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), aggregates, 0)
	})

	suite.Run("from과 to가 저장된 데이터를 포함하는 경우, 알맞은 개수의 데이터가 있어야 하고 에러가 없어야 한다", func() {
		// arrange
		start := time.Now()
		suite.storeStockAggregate(start, time.Second, 60)
		// act
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Second))
		defer cancel()

		aggregates, err := suite.testClient.FetchByTimeRange(ctx, testStockID, testTimeFrame, start, start.Add(30*time.Second))
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), aggregates, 30)
	})

	suite.Run("필수 필드가 빠졌을 때, 에러가 발생해야 한다.", func() {
		// arrange
		start := time.Now()
		suite.storeBrokenStockAggregate(start, time.Second, 60)
		// act
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Second))
		defer cancel()

		aggregates, err := suite.testClient.FetchByTimeRange(ctx, testStockID, "1m", start, start.Add(30*time.Second))
		assert.Error(suite.T(), err)
		assert.Len(suite.T(), aggregates, 0)
	})
}

func (suite *InfluxTestSuite) TestFetchLimitedTradeAfter() {
	suite.Run("빈 bucket에 쿼리를 요청한 경우, data가 비어있어야 하고 에러가 없어야 한다.", func() {
		// arrange
		// act
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
		now := time.Now()
		defer cancel()

		aggregates, err := suite.testClient.FetchLimitedTradeAfter(ctx, testStockID, "1m", now.Add(-5*60*time.Second), 10)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), aggregates, 0)
	})

	suite.Run("가져오려는 데이터보다 더 많은 데이터가 저장됐을 때, 원하는 만큼 데이터를 가져와야 한다.", func() {
		// arrange
		storeNum := 60
		storeInterval := time.Second
		start := time.Now().Add(time.Duration(-storeNum) * storeInterval)
		suite.storeStockAggregate(start, storeInterval, storeNum)

		// act
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Second))
		defer cancel()

		aggregates, err := suite.testClient.FetchLimitedTradeAfter(ctx, testStockID, testTimeFrame, start, 30)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), aggregates, 30)
	})

	suite.Run("가져오려는 데이터보다 더 적은 데이터가 저장됐을 때, 저장된 만큼 데이터를 가져와야 한다.", func() {
		// arrange
		storeNum := 20
		storeInterval := time.Second
		start := time.Now().Add(time.Duration(-storeNum) * storeInterval)
		suite.storeStockAggregate(start, storeInterval, storeNum)

		// act
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Second))
		defer cancel()

		aggregates, err := suite.testClient.FetchLimitedTradeAfter(ctx, testStockID, testTimeFrame, start, 30)
		assert.NoError(suite.T(), err)
		assert.Len(suite.T(), aggregates, 20)
	})

	suite.Run("필수 필드가 빠졌을 때, 에러가 발생해야 한다.", func() {
		// arrange
		storeNum := 60
		storeInterval := time.Second
		start := time.Now().Add(-time.Duration(storeNum) * storeInterval)
		suite.storeBrokenStockAggregate(start, storeInterval, storeNum)

		// act
		ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Second))
		defer cancel()

		aggregates, err := suite.testClient.FetchLimitedTradeAfter(ctx, testStockID, "1m", start, 30)
		assert.Error(suite.T(), err)
		assert.Len(suite.T(), aggregates, 0)
	})
}

func (suite *InfluxTestSuite) TearDownSuite() {
	suite.rawClient.Close()
}

func (suite *InfluxTestSuite) createBucketIfNotExits(orgName, bucketName string) error {
	org, err := suite.rawClient.OrganizationsAPI().FindOrganizationByName(context.Background(), orgName)
	if err != nil {
		return err
	}

	_, err = suite.rawClient.BucketsAPI().FindBucketByName(context.Background(), bucketName)

	if err != nil {
		_, err = suite.rawClient.BucketsAPI().CreateBucketWithName(context.Background(), org, bucketName)
		return err
	}
	return nil
}

func (suite *InfluxTestSuite) recreateBucket(orgName, bucketName string) error {
	org, err := suite.rawClient.OrganizationsAPI().FindOrganizationByName(context.Background(), orgName)
	if err != nil {
		return err
	}

	bucket, err := suite.rawClient.BucketsAPI().FindBucketByName(context.Background(), bucketName)
	if err != nil {
		return nil
	}

	suite.rawClient.BucketsAPI().DeleteBucket(context.Background(), bucket)
	_, err = suite.rawClient.BucketsAPI().CreateBucketWithName(context.Background(), org, bucketName)

	return err
}

func (suite *InfluxTestSuite) storeStockAggregate(start time.Time, interval time.Duration, num int) error {
	writer := suite.rawClient.WriteAPIBlocking(suite.options.Org, suite.options.TradeBucketName)
	for i := 0; i < num; i++ {
		err := writer.WritePoint(
			context.Background(),
			write.NewPoint(
				fmt.Sprintf("%s.%s", testStockID, testTimeFrame),
				map[string]string{},
				map[string]interface{}{
					"open":   float64(1.0),
					"close":  float64(2.0),
					"high":   float64(3.0),
					"low":    float64(4.0),
					"volume": float64(4),
				},
				start.Add(time.Duration(i)*interval),
			),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (suite *InfluxTestSuite) storeBrokenStockAggregate(start time.Time, interval time.Duration, num int) error {
	writer := suite.rawClient.WriteAPIBlocking(suite.options.Org, suite.options.TradeBucketName)
	for i := 0; i < num; i++ {
		err := writer.WritePoint(
			context.Background(),
			write.NewPoint(
				fmt.Sprintf("%s.%s", testStockID, testTimeFrame),
				map[string]string{},
				map[string]interface{}{
					"start": float64(1.0),
					"stop":  float64(2.0),
					"high":  float64(3.0),
					"low":   float64(4.0),
				},
				start.Add(time.Duration(i)*interval),
			),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func TestInfluxTestSuite(t *testing.T) {
	suite.Run(t, new(InfluxTestSuite))
}
