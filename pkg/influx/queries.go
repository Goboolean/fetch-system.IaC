package influx

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/Goboolean/fetch-system.IaC/pkg/model"
)

func (d *DB) FetchByTimeRange(
	ctx context.Context,
	productID string,
	timeFrame string,
	from time.Time,
	to time.Time) ([]*model.StockAggregate, error) {
	queryRes, err := d.reader.Query(ctx, fmt.Sprintf(
		`from(bucket:"%s")
				|> range(start:%d, stop:%d) 
				|> filter(fn: (r) => r._measurement == "%s.%s")
				|> filter(fn: (r) => (r._field == "open" or r._field == "close" or r._field == "high" or r._field == "low" or r._field == "volume"))
				|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")`,
		d.tradeBucket, from.Unix(), to.Unix(), productID, timeFrame))
	if err != nil {
		return nil, err
	}

	// When query result is empty but error is not occurred
	if queryRes == nil {
		return make([]*model.StockAggregate, 0), nil
	}

	data := []*model.StockAggregate{}
	for queryRes.Next() {
		aggregate := &model.StockAggregate{}

		if err := extractFieldValueByKey(queryRes.Record().Values(), "open", &aggregate.Open); err != nil {
			return nil, fmt.Errorf(`extracting field: can't extract "open": %w`, err)
		}
		if err := extractFieldValueByKey(queryRes.Record().Values(), "close", &aggregate.Close); err != nil {
			return nil, fmt.Errorf(`extracting field: can't extract "close": %v\n`, err)
		}
		if err := extractFieldValueByKey(queryRes.Record().Values(), "high", &aggregate.High); err != nil {
			return nil, fmt.Errorf(`extracting field: can't extract "high": %v\n`, err)
		}
		if err := extractFieldValueByKey(queryRes.Record().Values(), "low", &aggregate.Low); err != nil {
			return nil, fmt.Errorf(`extracting field: can't extract "low": %v\n`, err)
		}
		if err := extractFieldValueByKey(queryRes.Record().Values(), "volume", &aggregate.Volume); err != nil {
			return nil, fmt.Errorf(`extracting field: can't extract "volume": %v\n`, err)
		}
		if err := extractFieldValueByKey(queryRes.Record().Values(), "_time", &aggregate.Time); err != nil {
			return nil, fmt.Errorf(`extracting field: can't extract "_time": %v\n`, err)
		}

		data = append(data, aggregate)
	}

	return data, nil
}

// FetchLimitedTradeAfter fetches limited trade data for a specific product and time frame after a given start time.
// It queries the InfluxDB database to retrieve the stock aggregates (open, close, high, low, volume) for the specified product and time frame.
// The number of results is limited by the 'limit' parameter.
// The function returns a slice of StockAggregate pointers and an error if any occurred.
func (d *DB) FetchLimitedTradeAfter(
	ctx context.Context,
	productID string,
	timeFrame string,
	start time.Time,
	limit int) ([]*model.StockAggregate, error) {

	reader := d.client.QueryAPI("org")
	queryRes, err := reader.Query(ctx, fmt.Sprintf(
		`from(bucket:"%s")
			|> range(start:%d) 
			|> filter(fn: (r) => r._measurement == "%s.%s")
			|> filter(fn: (r) => (r._field == "open" or r._field == "close" or r._field == "high" or r._field == "low" or r._field == "volume"))
			|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")
			|> limit(n:%d)`,
		d.tradeBucket, start.Unix(), productID, timeFrame, limit))
	if err != nil {
		return nil, err
	}

	// When query result is empty but error is not occurred
	if queryRes == nil {
		return make([]*model.StockAggregate, 0), nil
	}

	data := []*model.StockAggregate{}
	for queryRes.Next() {
		aggregate := &model.StockAggregate{}

		if err := extractFieldValueByKey(queryRes.Record().Values(), "open", &aggregate.Open); err != nil {
			return nil, fmt.Errorf(`extracting field: can't extract "open": %w`, err)
		}
		if err := extractFieldValueByKey(queryRes.Record().Values(), "close", &aggregate.Close); err != nil {
			return nil, fmt.Errorf(`extracting field: can't extract "close": %v\n`, err)
		}
		if err := extractFieldValueByKey(queryRes.Record().Values(), "high", &aggregate.High); err != nil {
			return nil, fmt.Errorf(`extracting field: can't extract "high": %v\n`, err)
		}
		if err := extractFieldValueByKey(queryRes.Record().Values(), "low", &aggregate.Low); err != nil {
			return nil, fmt.Errorf(`extracting field: can't extract "low": %v\n`, err)
		}
		if err := extractFieldValueByKey(queryRes.Record().Values(), "volume", &aggregate.Volume); err != nil {
			return nil, fmt.Errorf(`extracting field: can't extract "volume": %v\n`, err)
		}
		if err := extractFieldValueByKey(queryRes.Record().Values(), "_time", &aggregate.Time); err != nil {
			return nil, fmt.Errorf(`extracting field: can't extract "_time": %v\n`, err)
		}

		data = append(data, aggregate)
	}

	return data, nil

}

func extractFieldValueByKey(values map[string]interface{}, key string, target any) error {

	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.IsNil() {
		return ErrInvalidFieldType
	}

	defer func() {
		if r := recover(); r != nil {
			targetValue.Elem().Set(reflect.Zero(targetValue.Elem().Type()))
		}
	}()

	val, ok := values[key]
	if !ok {
		return ErrFieldDoesNotExist
	}

	valueToSet := reflect.ValueOf(val)
	if !valueToSet.Type().AssignableTo(targetValue.Elem().Type()) {
		return ErrInvalidFieldType
	}

	targetValue.Elem().Set(valueToSet)
	return nil
}
