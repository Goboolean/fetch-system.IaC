package mongo_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.infrastructure/pkg/mongo"
	"github.com/stretchr/testify/assert"
	_ "github.com/Goboolean/common/pkg/env"
)



var db *mongo.DB

func SetupMongo() *mongo.DB {
	db, err := mongo.NewDB(&resolver.ConfigMap{
		"HOST":     os.Getenv("MONGO_HOST"),
		"USER":     os.Getenv("MONGO_USER"),
		"PORT":     os.Getenv("MONGO_PORT"),
		"PASSWORD": os.Getenv("MONGO_PASS"),
		"DATABASE": os.Getenv("MONGO_DATABASE"),
	})
	if err != nil {
		panic(err)
	}
	return db
}

func TeardownMongo(db *mongo.DB) {
	db.Close()
}

func TestMain(m *testing.M) {
	db = SetupMongo()
	code := m.Run()
	TeardownMongo(db)
	os.Exit(code)
}



func TestConstructor(t *testing.T) {
	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
		defer cancel()

		err := db.Ping(ctx)
		assert.NoError(t, err)
	})
}


func TestQueries(t *testing.T) {

	const productId = "stock.test.kor"
	const productType = "1m"

	t.Run("FetchAllStockBatch", func(t *testing.T) {
		ctx := context.Background()
		data, err := db.FetchAllStockBatch(ctx, productId, productType)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(data))
	})

	t.Run("FetchAllStockBatchMassive", func(t *testing.T) {
		var ch = make(chan *mongo.Aggregate, 0)

		ctx := context.Background()
		err := db.FetchAllStockBatchMassive(ctx, productId, productType, ch)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(ch))
	})
}