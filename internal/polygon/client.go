package polygon

import (
	"context"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)




type Client struct {
	conn *polygon.Client
}


func New(c *resolver.ConfigMap) (*Client, error) {

	key, err := c.GetStringKey("SECRET_KEY")
	if err != nil {
		return nil, err
	}

	conn := polygon.New(key)

	return &Client{
		conn: conn,
	}, nil
}


func (c *Client) GetAllProducts(ctx context.Context) ([]string, error){
	var includeOTC = true
	resp, err := c.conn.GetAllTickersSnapshot(ctx, &models.GetAllTickersSnapshotParams{
		Locale: models.US,
		MarketType: models.Stocks,
		IncludeOTC: &includeOTC,
	})
	if err != nil {
		return nil, err
	}

	productList := make([]string, len(resp.Tickers))

	for i, ticker := range resp.Tickers {
		productList[i] = ticker.Ticker
	}
	return productList, nil
}