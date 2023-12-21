package kis

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/Goboolean/common/pkg/resolver"
	"github.com/Goboolean/fetch-system.IaC/internal/model"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)




type Reader struct {
	filepath string
}

func New(c *resolver.ConfigMap) (*Reader, error) {
	filepath, err := c.GetStringKey("FILEPATH")
	if err != nil {
		return nil, err
	}
	return &Reader{
		filepath: filepath,
	}, nil
}



func (r *Reader) ReadAllTickerDetalis() ([]*model.TickerDetail, error) {

	var tickerDetails []*model.TickerDetail

	file, err := os.Open(r.filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := korean.EUCKR.NewDecoder()
	reader := csv.NewReader(transform.NewReader(bufio.NewReader(file), decoder))

	_, err = reader.Read()
	if err != nil {
		return nil, err
	}

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }

		tickerDetails = append(tickerDetails, &model.TickerDetail{
			Ticker: record[0],
			Name: record[1],
			Exchange: record[2],
		})
    }

	return tickerDetails, nil
}