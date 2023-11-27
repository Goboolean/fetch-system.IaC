package mongo_test

import (
	"context"
	"os"
	"testing"
	"time"

	_ "github.com/Goboolean/common/pkg/env"
	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.infrastructure/pkg/mongo"
	"github.com/stretchr/testify/assert"
)



var db *mongo.DB

func SetupMongo() *mongo.DB {
	db, err := mongo.NewDB(&resolver.ConfigMap{
		"HOST":     os.Getenv("MONGODB_HOST"),
		"PORT":     os.Getenv("MONGODB_PORT"),
		"USER":     os.Getenv("MONGODB_USERNAME"),
		"PASSWORD": os.Getenv("MONGODB_PASSWORD"),
		"DATABASE": os.Getenv("MONGODB_DATABASE"),
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

	const productId = "test.goboolean.kor"
	const productType = "1s"

	var data = &mongo.Aggregate{
		Timestamp: time.Now().UnixNano(),
	}

	t.Run("InsertProduct", func(t *testing.T) {
		ctx := context.Background()
		err := db.InsertProduct(ctx, productId, productType, data)
		assert.NoError(t, err)
	})

	t.Run("FetchAllStockBatch", func(t *testing.T) {
		ctx := context.Background()
		data, err := db.FetchAllStockBatch(ctx, productId, productType)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(data))
		t.Logf("data: %d", len(data))
	})

	t.Run("FetchAllStockBatchMassive", func(t *testing.T) {
		var ch = make(chan *mongo.Aggregate, 100)

		ctx := context.Background()
		err := db.FetchAllStockBatchMassive(ctx, productId, productType, ch)
		assert.NoError(t, err)
		assert.NotEqual(t, 0, len(ch))
		t.Logf("data: %d", len(ch))
	})
}