package etcd_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/Goboolean/fetch-system.IaC/internal/etcd"
	etcdutil "github.com/Goboolean/fetch-system.IaC/internal/etcd/util"
	"github.com/stretchr/testify/assert"
)



const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomString(n int) string {
    const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

    // Seed the random number generator
    rand.Seed(time.Now().UnixNano())

    b := make([]byte, n)
    for i := range b {
        b[i] = letterBytes[rand.Intn(len(letterBytes))]
    }
    return string(b)
}



func generateProduct() *etcd.Product {
	return &etcd.Product{
		ID:       randomString(32),
		Symbol:   randomString(16),
		Platform: randomString(8),
		Locale:   randomString(4),
		Market:   randomString(8),
	}
}

func generateProducts(n int) []*etcd.Product {
	products := make([]*etcd.Product, n)
	for i := 0; i < n; i++ {
		products[i] = generateProduct()
	}
	return products
}



func TestLoadDataOnInsertProductQuery(t *testing.T) {

	const count = 25 // maximum batch size
	var id string

	t.Run("ProductSize", func(t *testing.T) {
		product := generateProduct()
		data, err := etcdutil.Serialize(product)
		assert.NoError(t, err)

		var size = 0
		for k, v := range data {
			size += (len(v) + len(k))
		}
		t.Log("Product size:", size)
	})

	t.Run("InsertProducts", func(t *testing.T) {
		products := generateProducts(count)
		id = products[0].ID

		err := client.InsertProducts(context.Background(), products)
		assert.NoError(t, err)
	})

	t.Run("MeasureQueryTime: LowLoad", func(t *testing.T) {
		start := time.Now()

		_, err := client.GetProduct(context.Background(), id)
		assert.NoError(t, err)

		elasped := time.Since(start)
		t.Log("Elapsed time:", elasped)
	})

	t.Run("InsertMasiveProducts", func(t *testing.T) {
		products := generateProducts(10000)
		start := time.Now()

		for i := 0; i < len(products); i += count {
			j := i + count
			if j > len(products) {
				j = len(products)
			}
			err := client.InsertProducts(context.Background(), products[i:j])
			assert.NoError(t, err)
		}

		elasped := time.Since(start)
		t.Log("Elapsed time:", elasped)
	})

	t.Run("UpsertMasiveProducts", func(t *testing.T) {
		products := generateProducts(10000)
		start := time.Now()

		for i := 0; i < len(products); i += count {
			j := i + count
			if j > len(products) {
				j = len(products)
			}
			err := client.UpsertProducts(context.Background(), products[i:j])
			assert.NoError(t, err)
		}

		elasped := time.Since(start)
		t.Log("Elapsed time:", elasped)
	})



	t.Run("MeasureQueryTime: GetOne", func(t *testing.T) {
		start := time.Now()

		_, err := client.GetProduct(context.Background(), id)
		assert.NoError(t, err)

		elasped := time.Since(start)
		t.Log("Elapsed time:", elasped)
	})

	t.Run("MeasuerQueryTime: GetAll", func(t *testing.T) {
		start := time.Now()

		_, err := client.GetAllProducts(context.Background())
		assert.NoError(t, err)

		elasped := time.Since(start)
		t.Log("Elapsed time:", elasped)
	})

	t.Cleanup(func() {
		err := client.DeleteAllProducts(context.Background())
		assert.NoError(t, err)
	})
}



