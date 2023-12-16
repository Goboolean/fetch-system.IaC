package polygon

import (
	"context"
	"sync"

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


func (c *Client) GetTickerDetail(ctx context.Context, ticker string) (*TickerDetail, error) {

	resp, err := c.conn.GetTickerDetails(ctx, &models.GetTickerDetailsParams{
		Ticker: ticker,
	})
	if err != nil {
		return nil, err
	}

	return &TickerDetail{
		Ticker: resp.Results.Ticker,
		Name: resp.Results.Name,
		Description: resp.Results.Description,
		Exchange: resp.Results.PrimaryExchange,
	}, nil
}


var defaultSemaphoreSize = 100

func (c *Client) GetTickerDetailsMany(ctx context.Context, tickers []string) ([]*TickerDetailResult, error) {

	details := make([]*TickerDetailResult, len(tickers))

	semaphore := make(chan struct{}, defaultSemaphoreSize)
	wg := sync.WaitGroup{}

	for i, ticker := range tickers {
		semaphore <- struct{}{}
		wg.Add(1)

		go func(i int, ticker string) {
			defer func() {
				<-semaphore
				wg.Done()
			}()

			resp, err := c.conn.GetTickerDetails(ctx, &models.GetTickerDetailsParams{
				Ticker: ticker,
			})
			if err != nil {
				details[i] = &TickerDetailResult{
					Message: resp.ErrorMessage,
					Status: resp.Status,
				}
				return
			}

			details[i] = &TickerDetailResult{
				TickerDetail: TickerDetail{
					Ticker: resp.Results.Ticker,
					Name: resp.Results.Name,
					Description: resp.Results.Description,
					Exchange: resp.Results.PrimaryExchange,
				},
				Status: resp.Status,
				Message: resp.Message,
			}

		}(i, ticker)
	}

	wg.Wait()
	return details, nil
}