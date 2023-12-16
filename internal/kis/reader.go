package kis

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"

	"github.com/Goboolean/fetch-system.IaC/internal/model"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/transform"
)




type Reader struct {}



func (r *Reader) ReadAllTickerDetalis(filepath string) ([]*model.TickerDetail, error) {

	var tickerDetails []*model.TickerDetail

	file, err := os.Open(filepath)
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