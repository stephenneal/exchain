package microservice

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func makeRefreshTickerEndpoint(svc TickerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(refreshTickerRequest)
		svc.RefreshTicker(req.S)
		return emptyResponse{}, nil
	}
}

func makePrintTickersEndpoint(svc TickerService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		svc.PrintTickers()
		return emptyResponse{}, nil
	}
}

func decodeEmptyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request emptyRequest
	return request, nil
}

func decodeRefreshTickerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request refreshTickerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if _, ok := response.(emptyResponse); ok {
		return nil;
	}
	return json.NewEncoder(w).Encode(response)
}

type emptyRequest struct {
}

type emptyResponse struct {
}

type refreshTickerRequest struct {
	S string `json:"s"`
}
