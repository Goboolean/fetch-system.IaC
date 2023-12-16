package polygon_test

import (
	"context"
	"os"
	"testing"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/polygon"
	"github.com/stretchr/testify/assert"

	_ "github.com/Goboolean/common/pkg/env"
)




func SetupPolygon() *polygon.Client {

	p, err := polygon.New(&resolver.ConfigMap{
		"SECRET_KEY": os.Getenv("POLYGON_SECRET_KEY"),
	})
	if err != nil {
		panic(err)
	}

	return p
}



func TestMethod(t *testing.T) {

	p := SetupPolygon()

	var tickerList []string
	var tickerDetails []*polygon.TickerDetailResult

	t.Run("GetAllProducts", func(t *testing.T) {
		var err error
		tickerList, err = p.GetAllProducts(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, tickerList)
	})

	t.Run("GetTickerDetail", func(t *testing.T) {
		tickerDetail, err := p.GetTickerDetail(context.Background(), tickerList[0])
		assert.NoError(t, err)
		assert.NotEmpty(t, tickerDetail)
	})

	t.Run("GetTickerDetailMany", func(t *testing.T) {
		var err error
		tickerDetails, err = p.GetTickerDetailsMany(context.Background(), tickerList[:100])
		assert.NoError(t, err)
		assert.NotEmpty(t, tickerDetails)
	})
}