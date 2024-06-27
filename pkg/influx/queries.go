package influx

import (
	"context"
	"fmt"
	"time"

	"github.com/Goboolean/fetch-system.IaC/pkg/influx/mapper"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func (d *DB) FetchByTimeRange(
	ctx context.Context,
	productID string,
	timeFrame string,
	from time.Time,
	to time.Time) ([]map[string]interface{}, error) {
	queryRes, err := d.reader.Query(ctx, fmt.Sprintf(
		`from(bucket:"%s")
			|> range(start:%d, end:%d) 
			|> filter(fn: (r) => r._measurement == "%s.%s")
			|> filter(fn: (r) => (r._field == "open" or r._field == "close" or r._field == "high" or r._field == "low"))
			|> pivot(rowKey:["_time"], columnKey: ["_field"], valueColumn: "_value")`,
		tradeBucket, from.Unix(), to.Unix(), productID, timeFrame))
	if err != nil {
		return nil, err
	}

	data := []map[string]interface{}{}
	for queryRes.Next() {
		data = append(data, queryRes.Record().Values())
	}
	return data, nil
}

func (d *DB) GetOrderEventWriter(taskID string, event OrderEvent) error {

	d.orderWriter.WritePoint(write.NewPoint(
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

	return nil
}

func (d *DB) InsertAnnotation(taskID string, annotation any, createdAt time.Time) error {

	f, err := mapper.StructToPoint(annotation)
	if err != nil {
		return err
	}

	d.annotationWriter.WritePoint(write.NewPoint(
		taskID,
		map[string]string{},
		f,
		createdAt,
	))
	return nil
}
