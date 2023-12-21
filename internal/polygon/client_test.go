package polygon_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/model"
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

func TestClient(t *testing.T) {

	p := SetupPolygon()

	t.Run("Ping", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err := p.Ping(ctx)
		assert.NoError(t, err)
	})
}



func TestMethod(t *testing.T) {

	p := SetupPolygon()

	var tickerList []string
	var tickerDetails []*model.TickerDetailResult

	t.Run("GetAllTickers", func(t *testing.T) {
		var err error
		tickerList, err = p.GetAllTickers(context.Background())
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