package microservice

import (
	"context"
	"encoding/json"
	"net/http"

    "github.com/stephenneal/exchain/exchange"

	"github.com/go-kit/kit/endpoint"
)

func makeGetTickerEndpoint(svc TickerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(tickerRequest)
		err, tickers := svc.GetTicker(req.Pair)
		if (err != nil) {
		    return nil, err
		}
		return tickersResponse{ tickers }, nil
	}
}

func makeGetTickersEndpoint(svc TickerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		err, tickers := svc.GetTickers()
		if (err != nil) {
		    return nil, err
		}
		return tickersResponse{ tickers }, nil
	}
}

func decodeEmptyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request emptyRequest
	return request, nil
}

func decodeTickerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request tickerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

type emptyRequest struct {
}

type emptyResponse struct {
}

type tickerRequest struct {
	Pair string `json:"pair"`
}

type tickersResponse struct {
  	Tickers []exchange.Ticker `json:"tickers"`
}
