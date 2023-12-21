package kis_test

import (
	"testing"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/kis"
	"github.com/stretchr/testify/assert"
)


func SetupReader() *kis.Reader {
	r, err := kis.New(&resolver.ConfigMap{
		"FILEPATH": "../../api/csv/data.csv",
	})
	if err != nil {
		panic(err)
	}

	return r
}


func TestMethod(t *testing.T) {

	r := SetupReader()

	t.Run("ReadAllTickerDetalis", func(t *testing.T) {
		tickers, err := r.ReadAllTickerDetalis()
		assert.NoError(t, err)
		assert.NotEmpty(t, tickers)
		assert.NotEmpty(t, tickers[0].Ticker)
		assert.NotEmpty(t, tickers[0].Name)
		assert.NotEmpty(t, tickers[0].Exchange)
	})
}