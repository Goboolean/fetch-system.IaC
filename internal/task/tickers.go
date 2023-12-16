package task

import (
	"context"

	"github.com/Goboolean/fetch-system.IaC/internal/model"
)


func (m *Manager) filterTickerDetailsRespOK(resp []*model.TickerDetailResult) []*model.TickerDetail {
	var tickerDetails []*model.TickerDetail
	for _, r := range resp {
		if r.Status == "OK" {
			tickerDetails = append(tickerDetails, &r.TickerDetail)
		}
	}
	return tickerDetails
}

func (m *Manager) filterTickersRespError(resp []*model.TickerDetailResult) []string {
	var tickerDetails []string
	for _, r := range resp {
		if r.Status != "OK" {
			tickerDetails = append(tickerDetails, r.TickerDetail.Ticker)
		}
	}
	return tickerDetails
}


const retryCount = 3

func (m *Manager) GetAllUSATickerDetails(ctx context.Context) ([]*model.TickerDetail, error) {
	var err error

	var tickerDetails []*model.TickerDetail
	var tickers []string

	tickers, err = m.polygon.GetAllTickers(ctx)
	if err != nil {
		return nil, err
	}

	for i := 0; i < retryCount; i++ {
		resp, err := m.polygon.GetTickerDetailsMany(ctx, tickers)
		if err != nil {
			return nil, err
		}
	
		tickerDetails = append(tickerDetails, m.filterTickerDetailsRespOK(resp)...)

		tickers = m.filterTickersRespError(resp)

		if len(tickers) == 0 {
			break
		}
	}

	return tickerDetails, nil
}