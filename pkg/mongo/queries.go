package mongo

import (
	"context"
	"fmt"

	"github.com/Goboolean/fetch-system.IaC/api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (db *DB) FetchAllStockBatch(ctx context.Context, productId string, timeFrame string) ([]*Aggregate, error) {
	session, err := db.client.StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(ctx)

	symbol := fmt.Sprintf("%s.%s", productId, timeFrame)
	if valid := model.IsSymbolValid(symbol); !valid {
		return nil, ErrInvalidSymbol
	}
	coll := db.client.Database(db.defaultDB).Collection(symbol)

	results := make([]*Aggregate, 0)

	return results, mongo.WithSession(ctx, session, func(ctx mongo.SessionContext) error {
		cursor, err := coll.Find(ctx, bson.M{})
		if err != nil {
			return err
		}

		for cursor.Next(ctx) {
			var data Aggregate
			if err := cursor.Decode(&data); err != nil {
				return err
			}

			results = append(results, &data)
		}

		return cursor.Close(ctx)
	})
}

func (db *DB) FetchAllStockBatchMassive(ctx context.Context, productId string, timeFrame string, productCh chan<- *Aggregate) error {
	session, err := db.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	symbol := fmt.Sprintf("%s.%s", productId, timeFrame)
	if valid := model.IsSymbolValid(symbol); !valid {
		return ErrInvalidSymbol
	}
	coll := db.client.Database(db.defaultDB).Collection(symbol)

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

func (db *DB) InsertProduct(ctx context.Context, productId string, timeFrame string, data *Aggregate) error {
	session, err := db.client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	symbol := fmt.Sprintf("%s.%s", productId, timeFrame)
	if valid := model.IsSymbolValid(symbol); !valid {
		return ErrInvalidSymbol
	}

	coll := db.client.Database(db.defaultDB).Collection(symbol)

	return mongo.WithSession(ctx, session, func(ctx mongo.SessionContext) error {
		_, err := coll.InsertOne(ctx, data)
		return err
	})
}
