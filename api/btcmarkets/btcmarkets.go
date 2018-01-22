package btcmarkets

import (
    "fmt"

    "github.com/romana/rlog"

    "github.com/stephenneal/exchain/api"
)

const (
    BASE_URL = "https://api.btcmarkets.net"
)

type TickerResponse struct {
	BestBid    float64 `json:"bestBid"`
	BestAsk    float64 `json:"bestAsk"`
	LastPrice  float64 `json:"lastPrice"`
	Currency   string  `json:"currency"`
	Instrument string  `json:"instrument"`
	Timestamp  int     `json:"timestamp"`
	Volume24H  float64 `json:"volume24h"`
}

func Ticker(pair string) {
    var response TickerResponse

    if err := api.GetJson(BASE_URL + "/market/" + pair + "/tick", &response); err != nil {
        rlog.Error(err)
    } else {
        rlog.Info(fmt.Sprintf("BTCMarkets (%s); Last=%f", pair, response.LastPrice))
    }
}
