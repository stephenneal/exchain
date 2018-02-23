package exchangems

import (
    "context"

    "github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
    TickersEndpoint     endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
    return Endpoints{
        TickersEndpoint:    MakeTickersEndpoint(s),
    }
}

func MakeTickersEndpoint(svc Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        req := request.(tickerRequest)
        err, tickers := svc.GetTickers(req.Base, req.Quot)
        return tickersResponse{Tickers: tickers, Err: err}, nil
    }
}

/*
func MakeGetTickersEndpoint(svc Service) endpoint.Endpoint {
    return func(ctx context.Context, request interface{}) (interface{}, error) {
        err, tickers := svc.GetTickers()
        return tickerSummaryResponse{Tickers: tickers, Err: err}, nil
    }
}
*/

type emptyRequest struct {
}

type emptyResponse struct {
}

type tickerRequest struct {
    Base string
    Quot string
}

type tickersResponse struct {
    Tickers []Ticker `json:"tickers"`
    Err     error  `json:"error,omitempty"`
}

/*
type tickerSummaryResponse struct {
    Tickers []TickerSummary `json:"tickers"`
    Err     error  `json:"error,omitempty"`
}
*/

func (r tickersResponse) error() error { return r.Err }

//func (r tickerSummaryResponse) error() error { return r.Err }
