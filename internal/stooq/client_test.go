package stooq

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

func TestClient_GetStockData(t *testing.T) {
	type fields struct {
		Client http.Client
	}

	type args struct {
		ctx       context.Context
		stockCode string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Stock
		wantErr bool
		err     error
	}{
		{"success",
			fields{http.Client{}},
			args{context.Background(), "aapl.us"},
			&Stock{
				Symbol: "AAPL.US",
				Date:   "2022-10-07",
				Time:   "22:00:06",
				Open:   "142.54",
				High:   "143.15",
				Low:    "139.44",
				Close:  "140.09",
				Volume: "85925559",
			},
			false,
			nil},
		{"not-found",
			fields{http.Client{}},
			args{context.Background(), "aapl.uk"},
			nil,
			true,
			errors.New("error getting stock data: 404")},
		{"not-found-N/D",
			fields{http.Client{}},
			args{context.Background(), "aapl.ul"},
			&Stock{
				Symbol: "AAPL.UL",
				Date:   "N/D",
				Time:   "N/D",
				Open:   "N/D",
				High:   "N/D",
				Low:    "N/D",
				Close:  "N/D",
				Volume: "N/D",
			},
			true,
			nil,
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", "https://stooq.com/q/l/?s=aapl.us&f=sd2t2ohlcv&h&e=csv",
		httpmock.NewStringResponder(200, `Symbol,Date,Time,Open,High,Low,Close,Volume
AAPL.US,2022-10-07,22:00:06,142.54,143.15,139.44,140.09,85925559`))
	httpmock.RegisterResponder("GET", "https://stooq.com/q/l/?s=aapl.uk&f=sd2t2ohlcv&h&e=csv",
		httpmock.NewStringResponder(404, ""))
	httpmock.RegisterResponder("GET", "https://stooq.com/q/l/?s=aapl.ul&f=sd2t2ohlcv&h&e=csv",
		httpmock.NewStringResponder(200, `Symbol,Date,Time,Open,High,Low,Close,Volume
AAPL.UL,N/D,N/D,N/D,N/D,N/D,N/D,N/D`))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{
				Client: tt.fields.Client,
			}

			got, err := c.GetStockData(tt.args.ctx, tt.args.stockCode)
			if (err != nil) != tt.wantErr {
				assert.Equal(t, tt.err, err)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}
