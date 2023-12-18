package retrieve

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


func TestTickers(t *testing.T) {

	var count int

	m := new(Manager)
	m.polygon = SetupPolygon()
	
	t.Run("GetAllTickers", func(t *testing.T) {
		p := SetupPolygon()
		tickers, err := p.GetAllTickers(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, tickers)
		count = len(tickers)
	})

	t.Run("GetAllUSATickerDetails", func(t *testing.T) {
		t.Skip("Skip this test since it generates too much traffic to the API")

		tickerDetails, err := m.GetAllUSATickerDetails(context.Background())
		assert.NoError(t, err)
		assert.NotEmpty(t, tickerDetails)

		assert.GreaterOrEqual(t, count, len(tickerDetails))

		for _, tickerDetail := range tickerDetails {
			assert.NotEmpty(t, tickerDetail.Ticker)
			assert.NotEmpty(t, tickerDetail.Name)
		}
	})
}