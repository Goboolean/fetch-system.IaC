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
	"github.com/stretchr/testify/suite"
)

var testStockID = "stock.aapl.usa"
var testTimeFrame = "1m"

type DBTestSuite struct {
	suite.Suite
	testClient *influx.DB
	rawClient  influxdb2.Client
	options    influx.Opts
}

func (suite *DBTestSuite) SetupSuite() {
	suite.options = influx.Opts{
		URL:             os.Getenv("INFLUXDB_URL"),
		Token:           os.Getenv("INFLUXDB_TOKEN"),
		TradeBucketName: os.Getenv("INFLUXDB_TRADE_BUCKET"),
		Org:             os.Getenv("INFLUXDB_ORG"),
	}

	suite.rawClient = influxdb2.NewClient(suite.options.URL, suite.options.Token)

	var err error
	suite.testClient, err = influx.NewDB(&suite.options)
	suite.Require().NoError(err)

	err = suite.createBucketIfNotExits(suite.options.Org, suite.options.TradeBucketName)
	suite.Require().NoError(err)
}

func (suite *DBTestSuite) SetupTest() {
	err := suite.recreateBucket(suite.options.Org, suite.options.TradeBucketName)
	suite.Require().NoError(err)
}
func (suite *DBTestSuite) TestConstructor() {
	err := suite.testClient.Ping(context.Background())
	suite.NoError(err)
}

func (suite *DBTestSuite) TestFetchByTimeRange_ShouldReturnEmptyTradeWithoutError_WhenBucketIsEmpty() {
	// arrange
	// act
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	now := time.Now()
	defer cancel()

	aggregates, err := suite.testClient.FetchByTimeRange(ctx, testStockID, "1m", now.Add(-5*60*time.Second), now)
	suite.NoError(err)
	suite.Len(aggregates, 0)
}

func (suite *DBTestSuite) TestFetchByTimeRange_ShouldReturnEmptyTrade_WhenFromAndToRepresentOutOfStoredRange() {
	// arrange
	start := time.Now()
	err := suite.storeStockAggregate(start, time.Second, 60)
	suite.Require().NoError(err)
	// act
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	defer cancel()

	aggregates, err := suite.testClient.FetchByTimeRange(ctx, testStockID, testTimeFrame, start.Add(-5*time.Second), start.Add(-1*time.Second))
	suite.NoError(err)
	suite.Len(aggregates, 0)
	fmt.Println("end: ", "TestFetchByTimeRange_ShouldReturnEmptyTrade_WhenFromAndToRepresentOutOfStoredRange")
}

func (suite *DBTestSuite) TestFetchByTimeRange_ShouldFetchData_WhenFromAndToIncludeRangeOfStoredData() {
	// arrange
	start := time.Now()
	err := suite.storeStockAggregate(start, time.Second, 60)
	suite.Require().NoError(err)
	// act
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Second))
	defer cancel()

	aggregates, err := suite.testClient.FetchByTimeRange(ctx, testStockID, testTimeFrame, start, start.Add(30*time.Second))
	suite.NoError(err)
	suite.Len(aggregates, 30)
}
func (suite *DBTestSuite) TestFetchByTimeRange_ShouldReturnError_WhenRequiredFieldNotExists() {
	// arrange
	start := time.Now()
	err := suite.storeBrokenStockAggregate(start, time.Second, 60)
	suite.Require().NoError(err)
	// act
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Second))
	defer cancel()

	aggregates, err := suite.testClient.FetchByTimeRange(ctx, testStockID, "1m", start, start.Add(30*time.Second))
	suite.Error(err)
	suite.Len(aggregates, 0)
}

func (suite *DBTestSuite) TestFetchLimitedTradeAfter_ShouldReturnEmptyTradeWithoutError_WhenBucketIsEmpty() {
	// arrange
	// act
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
	now := time.Now()
	defer cancel()

	aggregates, err := suite.testClient.FetchLimitedTradeAfter(ctx, testStockID, "1m", now.Add(-5*60*time.Second), 10)
	suite.NoError(err)
	suite.Len(aggregates, 0)
}

func (suite *DBTestSuite) TestFetchLimitedTradeAfter_ShouldFetchLimitedAmountOfData_WhenMoreDataIsStored() {
	// arrange
	storeNum := 60
	storeInterval := time.Second
	start := time.Now().Add(time.Duration(-storeNum) * storeInterval)
	err := suite.storeStockAggregate(start, storeInterval, storeNum)
	suite.Require().NoError(err)

	// act
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Second))
	defer cancel()

	aggregates, err := suite.testClient.FetchLimitedTradeAfter(ctx, testStockID, testTimeFrame, start, 30)
	suite.NoError(err)
	suite.Len(aggregates, 30)
}

func (suite *DBTestSuite) TestFetchLimitedTradeAfter_ShouldFetchAllData_WhenRequestedAmountExceedsStoredData() {
	// arrange
	storeNum := 20
	storeInterval := time.Second
	start := time.Now().Add(time.Duration(-storeNum) * storeInterval)
	err := suite.storeStockAggregate(start, storeInterval, storeNum)
	suite.Require().NoError(err)

	// act
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Second))
	defer cancel()

	aggregates, err := suite.testClient.FetchLimitedTradeAfter(ctx, testStockID, testTimeFrame, start, 30)
	suite.NoError(err)
	suite.Len(aggregates, 20)
}

func (suite *DBTestSuite) TestFetchLimitedTradeAfter_ShouldFetchAllData_ShouldReturnError_WhenRequiredFieldNotExists() {
	// arrange
	storeNum := 60
	storeInterval := time.Second
	start := time.Now().Add(-time.Duration(storeNum) * storeInterval)
	err := suite.storeBrokenStockAggregate(start, storeInterval, storeNum)
	suite.Require().NoError(err)

	// act
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(100*time.Second))
	defer cancel()

	aggregates, err := suite.testClient.FetchLimitedTradeAfter(ctx, testStockID, "1m", start, 30)
	suite.Error(err)
	suite.Len(aggregates, 0)
}

func (suite *DBTestSuite) TearDownSuite() {
	suite.rawClient.Close()
}

func (suite *DBTestSuite) createBucketIfNotExits(orgName, bucketName string) error {
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

func (suite *DBTestSuite) recreateBucket(orgName, bucketName string) error {
	org, err := suite.rawClient.OrganizationsAPI().FindOrganizationByName(context.Background(), orgName)
	if err != nil {
		return err
	}

	bucket, err := suite.rawClient.BucketsAPI().FindBucketByName(context.Background(), bucketName)
	if err != nil {
		return nil
	}

	err = suite.rawClient.BucketsAPI().DeleteBucket(context.Background(), bucket)
	if err != nil {
		return nil
	}
	_, err = suite.rawClient.BucketsAPI().CreateBucketWithName(context.Background(), org, bucketName)

	return err
}

func (suite *DBTestSuite) storeStockAggregate(start time.Time, interval time.Duration, num int) error {
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

func (suite *DBTestSuite) storeBrokenStockAggregate(start time.Time, interval time.Duration, num int) error {
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
func TestInflux(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
