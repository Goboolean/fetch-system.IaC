package db_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/pkg/db"
	"github.com/stretchr/testify/assert"

	_ "github.com/Goboolean/common/pkg/env"
)



func SetupPostgreSQL() *db.Client {
	c, err := db.NewDB(&resolver.ConfigMap{
		"HOST": os.Getenv("POSTGRES_HOST"),
		"PORT": os.Getenv("POSTGRES_PORT"),
		"USER": os.Getenv("POSTGRES_USER"),
		"PASSWORD": os.Getenv("POSTGRES_PASSWORD"),
		"DATABASE": os.Getenv("POSTGRES_DATABASE"),
	})

	if err != nil {
		panic(err)
	}

	return c
}

func TeardownPostgreSQL(c *db.Client) {
	c.Close()
}



func TestClient(t *testing.T) {

	c := SetupPostgreSQL()
	defer TeardownPostgreSQL(c)

	t.Run("Ping()", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := c.Ping(ctx)
		assert.NoError(t, err)
	})
}




func TestInsertScenario(t *testing.T) {

	c := SetupPostgreSQL()

	t.Cleanup(func() {
		err := c.DeleteAllProducts(context.Background())
		assert.NoError(t, err)

		TeardownPostgreSQL(c)
	})

	var products = []db.InsertProductsParams{
		{
			ID: "stock.goboolean.test",
			Symbol: "goboolean",
			Locale: "test",
			Market: "stock",
		},
		{
			ID: "stock.golution.test",
			Symbol: "goboolean",
			Locale: "test",
			Market: "stock",
		},
	}

	t.Run("InsertProducts()", func(t *testing.T) {
		v, err := c.InsertProducts(context.Background(), products)
		assert.NoError(t, err)
		assert.Equal(t, int64(2), v)
	})

	t.Run("GetAllProducts()", func(t *testing.T) {
		results, err := c.GetAllProducts(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 2, len(results))
	})
}