package kis_test

import (
	"testing"

	"github.com/Goboolean/fetch-system.IaC/internal/kis"
	"github.com/stretchr/testify/assert"
)


func TestMethod(t *testing.T) {

	r := kis.Reader{}

	t.Run("ReadAllTickerDetalis", func(t *testing.T) {
		tickers, err := r.ReadAllTickerDetalis("../../api/csv/data.csv")
		assert.NoError(t, err)
		assert.NotEmpty(t, tickers)
		assert.NotEmpty(t, tickers[0].Ticker)
		assert.NotEmpty(t, tickers[0].Name)
		assert.NotEmpty(t, tickers[0].Exchange)
	})
}