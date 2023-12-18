package retrieve

import (
	"context"
	"fmt"

	"github.com/Goboolean/fetch-system.IaC/internal/kis"
	"github.com/Goboolean/fetch-system.IaC/internal/model"
	"github.com/Goboolean/fetch-system.IaC/internal/polygon"
	"github.com/Goboolean/fetch-system.IaC/pkg/db"
)






type Manager struct {
	polygon *polygon.Client
	db 	*db.Client
	kis *kis.Reader
}





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


func (m *Manager) InitUSAStocks(ctx context.Context) error {

	details, err := m.GetAllUSATickerDetails(ctx)
	if err != nil {
		return err
	}

	dtos := make([]db.InsertProductsParams, len(details))

	for i, detail := range details {
		dtos[i] = db.InsertProductsParams{
			ID:     fmt.Sprintf("%s.%s.%s", "stock", detail.Ticker, "usa"),
			Symbol: detail.Ticker,
			Locale: "usa",
			Market: "stock",
		}
	}

	count, err := m.db.InsertProducts(ctx, dtos)
	if err != nil {
		return err
	}

	if int(count) != len(details) {
		return fmt.Errorf("failed to insert products")
	}
	return nil
}


func (m *Manager) InitKORStocks(ctx context.Context) error {

	details, err := m.kis.ReadAllTickerDetalis()
	if err != nil {
		return err
	}

	dtos := make([]db.InsertProductsParams, len(details))

	for i, detail := range details {
		dtos[i] = db.InsertProductsParams{
			ID:     fmt.Sprintf("%s.%s.%s", "stock", detail.Ticker, "kor"),
			Symbol: detail.Ticker,
			Locale: "kor",
			Market: "stock",
		}
	}

	count, err := m.db.InsertProducts(ctx, dtos)
	if err != nil {
		return err
	}

	if int(count) != len(details) {
		return fmt.Errorf("failed to insert products")
	}
	return nil
}
