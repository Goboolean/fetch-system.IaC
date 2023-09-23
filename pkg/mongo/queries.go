package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)






func (db *DB) FetchAllStockBatch(ctx context.Context, productId string, productType string) ([]*Aggregate, error) {
	session, err := db.client.StartSession()
	if err != nil {
		return nil, err
	}

	collName := fmt.Sprintf("%s.%s", productId, productType)
	coll := db.client.Database(db.defaultDB).Collection(collName)

	results := make([]*Aggregate, 0)

	return results, mongo.WithSession(ctx, session, func(ctx mongo.SessionContext) error {
		cursor, err := coll.Find(ctx, bson.M{})
		if err != nil {
			return err
		}

		for cursor.Next(ctx) {
			var data Aggregate
			if err := cursor.Decode(data); err != nil {
				return err
			}

			results = append(results, &data)
		}

		return cursor.Close(ctx)
	})
}



func (db *DB) FetchAllStockBatchMassive(ctx context.Context, productId string, productType string, productCh chan<- *Aggregate) error {
	session, err := db.client.StartSession()
	if err != nil {
		return err
	}

	collName := fmt.Sprintf("%s.%s", productId, productType)
	coll := db.client.Database(db.defaultDB).Collection(collName)

	return mongo.WithSession(ctx, session, func(ctx mongo.SessionContext) error {
		cursor, err := coll.Find(ctx, bson.M{})
		if err != nil {
			return err
		}

		for cursor.Next(ctx) {
			var data Aggregate
			if err := cursor.Decode(&data); err != nil {
				return err
			}

			productCh <- &data
		}

		return cursor.Close(ctx)
	})
}