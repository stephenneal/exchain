package exchangems

import (
	"context"

    "github.com/stephenneal/exchain/exchange"

	"github.com/go-kit/kit/endpoint"
)

type emptyRequest struct {
}

type emptyResponse struct {
}

type tickerRequest struct {
	Pair string
}

type tickersResponse struct {
  	Tickers []exchange.Ticker `json:"tickers"`
	Err   error  `json:"error,omitempty"`
}

type tickerSummaryResponse struct {
  	Tickers []exchange.TickerSummary `json:"tickers"`
	Err   error  `json:"error,omitempty"`
}

func (r tickersResponse) error() error { return r.Err }

func makeGetTickerEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(tickerRequest)
		err, tickers := svc.GetTicker(req.Pair)
		return tickersResponse{Tickers: tickers, Err: err}, nil
	}
}

func makeGetTickersEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err, tickers := svc.GetTickers()
		return tickerSummaryResponse{Tickers: tickers, Err: err}, nil
	}
}
