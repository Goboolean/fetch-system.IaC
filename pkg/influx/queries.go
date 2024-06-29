package influx

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/Goboolean/fetch-system.IaC/pkg/model"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func (d *DB) FetchByTimeRange(
	ctx context.Context,
	productID string,
	timeFrame string,
	from time.Time,
	to time.Time) ([]*model.StockAggregate, error) {
	queryRes, err := d.reader.Query(ctx, fmt.Sprintf(
		`from(bucket:"%s")
			|> range(start:%d, end:%d) 
			|> filter(fn: (r) => r._measurement == "%s.%s")
			|> filter(fn: (r) => (r._field == "open" or r._field == "close" or r._field == "high" or r._field == "low"))
			|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")`,
		d.tradeBucket, from.Unix(), to.Unix(), productID, timeFrame))
	if err != nil {
		return nil, err
	}

	// When query result is empty but error is not occurred
	if queryRes.Record() == nil {
		return nil, nil
	}

	data := []*model.StockAggregate{}
	for queryRes.Next() {
		aggregate := &model.StockAggregate{}

		extractFieldValueByKey(queryRes.Record().Values(), "open", &aggregate.Open)
		extractFieldValueByKey(queryRes.Record().Values(), "close", &aggregate.Close)
		extractFieldValueByKey(queryRes.Record().Values(), "high", &aggregate.High)
		extractFieldValueByKey(queryRes.Record().Values(), "low", &aggregate.Low)
		extractFieldValueByKey(queryRes.Record().Values(), "_time", &aggregate.Time)
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

func (d *DB) InsertOrderEvent(ctx context.Context, taskID string, event OrderEvent) error {

	return d.orderWriter.WritePoint(ctx, write.NewPoint(
		taskID,
		map[string]string{},
		map[string]interface{}{
			"productID":         event.ProductID,
			"proportionPercent": event.ProportionPercent,
			"action":            event.Action,
			"task":              event.Task,
		},
		event.CreatedAt,
	))
}

func (d *DB) InsertAnnotation(ctx context.Context, taskID string, annotation map[string]interface{}, createdAt time.Time) error {

	return d.annotationWriter.WritePoint(ctx, write.NewPoint(
		taskID,
		map[string]string{},
		annotation,
		createdAt,
	))

}
