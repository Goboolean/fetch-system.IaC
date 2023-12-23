package dbiniter

import (
	"context"
	"fmt"

	"github.com/Goboolean/fetch-system.IaC/internal/kis"
	"github.com/Goboolean/fetch-system.IaC/internal/model"
	"github.com/Goboolean/fetch-system.IaC/internal/polygon"
	"github.com/Goboolean/fetch-system.IaC/pkg/db"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)






type Manager struct {
	polygon *polygon.Client
	db 	*db.Client
	kis *kis.Reader
}

func New(polygon *polygon.Client, db *db.Client, kis *kis.Reader) *Manager {
	return &Manager{
		polygon: polygon,
		db: db,
		kis: kis,
	}
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
		return nil, errors.Wrap(err, "Failed to get all tickers from polygon")
	}

	log.Infof("Received number of %d tickers", len(tickers))

	log.Info("Getting ticker details from polygon")
	for i := 0; i < retryCount; i++ {
		resp, err := m.polygon.GetTickerDetailsMany(ctx, tickers)
		if err != nil { 
			return nil, errors.Wrap(err, fmt.Sprintf("Failed to get ticker details from polygon (retry count: %d)", i))
		}
	
		tickerDetails = append(tickerDetails, m.filterTickerDetailsRespOK(resp)...)

		tickers = m.filterTickersRespError(resp)

		if len(tickers) == 0 {
			break
		}
	}
	log.Info("Getting ticker details from polygon finished")
	return tickerDetails, nil
}


func (m *Manager) CheckUSAStockStored(ctx context.Context) (bool, error) {

	count, err := m.db.CountProducts(ctx, db.CountProductsParams{
		Platform: db.PlatformPOLYGON,
		Market: db.MarketSTOCK,
	})
	if err != nil {
		return false, errors.Wrap(err, "Failed to count products")
	}

	if count == 0 {
		log.Info("USA stocks is not stored yet")
		return false, nil
	} else {
		log.Info("Already stored USA stocks")
		return true, nil
	}
}


func (m *Manager) StoreUSAStocks(ctx context.Context) error {

	details, err := m.GetAllUSATickerDetails(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get all ticker details from polygon")
	}

	dtos := make([]db.InsertProductsParams, len(details))

	for i, detail := range details {
		dtos[i] = db.InsertProductsParams{
			ID:     fmt.Sprintf("%s.%s.%s", "stock", detail.Ticker, "usa"),
			Symbol: detail.Ticker,
			Locale: db.LocaleUSA,
			Market: db.MarketSTOCK,
			Platform: db.PlatformPOLYGON,
		}
	}

	if _, err = m.db.InsertProducts(ctx, dtos); err != nil {
		return errors.Wrap(err, "Failed to insert products on database")
	}

	log.Infof("Successfully stored number of %d USA products", len(details))
	return nil
}


func (m *Manager) CheckKORStockStored(ctx context.Context) (bool, error) {

	count, err := m.db.CountProducts(ctx, db.CountProductsParams{
		Platform: db.PlatformKIS,
		Market: db.MarketSTOCK,
	})
	if err != nil {
		return false, errors.Wrap(err, "Failed to count products")
	}

	if count == 0 {
		log.Info("KOR stocks is not stored yet")
		return false, nil
	} else {
		log.Info("Already stored KOR stocks")
		return true, nil
	}
}


func (m *Manager) StoreKORStocks(ctx context.Context) error {

	details, err := m.kis.ReadAllTickerDetalis()
	if err != nil {
		return errors.Wrap(err, "Failed to read all ticker details from kis")
	}

	dtos := make([]db.InsertProductsParams, len(details))

	for i, detail := range details {
		dtos[i] = db.InsertProductsParams{
			ID:     fmt.Sprintf("%s.%s.%s", "stock", detail.Ticker, "kor"),
			Symbol: detail.Ticker,
			Locale: db.LocaleKOR,
			Market: db.MarketSTOCK,
			Platform: db.PlatformKIS,
			Name: pgtype.Text{
				String: detail.Name,
				Valid: (detail.Name != ""),
			},
			Description: pgtype.Text{
				String: detail.Description,
				Valid: (detail.Description != ""),
			},
		}
	}

	if _, err = m.db.InsertProducts(ctx, dtos); err != nil {
		return errors.Wrap(err, "Failed to insert products on database")
	}

	log.Infof("Successfully stored number of %d KOR products", len(details))
	return nil
}
