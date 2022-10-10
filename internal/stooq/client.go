package stooq

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gocarina/gocsv"
	log "github.com/sirupsen/logrus"
)

type Stock struct {
	Symbol string `csv:"Symbol"`
	Date   string `csv:"Date"`
	Time   string `csv:"Time"`
	Open   string `csv:"Open"`
	High   string `csv:"High"`
	Low    string `csv:"Low"`
	Close  string `csv:"Close"`
	Volume string `csv:"Volume"`
}

type Client struct {
	http.Client
}

func NewClient() *Client {
	return &Client{
		http.Client{},
	}
}

func (c *Client) GetStockData(ctx context.Context, stockCode string) (*Stock, error) {
	requestURL := fmt.Sprintf("https://stooq.com/q/l/?s=%s&f=sd2t2ohlcv&h&e=csv", stockCode)

	r, _ := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)

	response, err := c.Do(r)
	if err != nil {
		log.Errorf("error getting stock data: %s", err)

		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		log.Errorf("error getting stock data: %s", response.Status)

		return nil, fmt.Errorf("error getting stock data: %s", response.Status) //nolint:goerr113
	}

	defer func() {
		if err := response.Body.Close(); err != nil {
			log.Printf("error closing response body: %s\n", err)
		}
	}()

	var stocks []*Stock

	if err := gocsv.Unmarshal(response.Body, &stocks); err != nil {
		log.Errorf("error unmarshalling csv: %s", err)

		return nil, err
	}

	// we only need the first element, since we are only requesting one stock at a time
	return stocks[0], nil
}
