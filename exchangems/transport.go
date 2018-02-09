package exchangems

import (
	"context"
	"encoding/json"
	//"errors"
	"net/http"

	"github.com/gorilla/mux"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

//var errBadRoute = errors.New("bad route")

// MakeHandler returns a handler for the booking service.
func MakeHTTPHandler(s Service, logger kitlog.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)

	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

    r.Methods("GET").Path("/pub/v1/tickers/{base}").Handler(kithttp.NewServer(
        e.GetTickersEndpoint,
        decodeTickerRequest,
        kithttp.EncodeJSONResponse,
        opts...,
    ))
    /*
    getTickersHandler := kithttp.NewServer(
        makeGetTickersEndpoint(es),
        decodeEmptyRequest,
        kithttp.EncodeJSONResponse,
        opts...,
    )
    */

    //r.Handle("/pub/ex/v1/getTicker/{pair}", getTickerHandler).Methods("GET")
    //r.Handle("/pub/ex/v1/getTickers", getTickersHandler).Methods("GET")

	return r
}

func decodeEmptyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request emptyRequest
	return request, nil
}

func decodeTickerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	base, _ := vars["base"]
    quot, _ := vars["quot"]
	return tickerRequest{ Base: base, Quot: quot }, nil
    /*
	var request tickerRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
	*/
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

    w.WriteHeader(http.StatusInternalServerError)
    /* TODO...
	switch err {
	case cargo.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	case ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	*/
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}