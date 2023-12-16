package kis_test

import (
	"testing"

	"github.com/Goboolean/fetch-system.IaC/internal/kis"
	"github.com/stretchr/testify/assert"
)


func TestReadKORTickers(t *testing.T) {

	r := kis.Reader{}

	t.Run("ReadKORTickers", func(t *testing.T) {
		tickers, err := r.ReadAllTickerDetalis("../../api/csv/data.csv")
		assert.NoError(t, err)
		assert.NotEmpty(t, tickers)
	})
}