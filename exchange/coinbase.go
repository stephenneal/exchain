package exchange

import (
    "strings"
)

const (
    coinbaseName = "Coinbase"
)

var (
    coinbasePairs = []string{
        BTC_AUD,
        BTC_USD,
        ETH_AUD,
        ETH_USD,
    }
)

type coinbaseService struct{}

type coinbaseTicker struct {
    Data struct {
        Base       string `json:"base"`
        Currency   string `json:"currency"`
        Last       float64 `json:"amount,string"`
    } `json:"data"`
}

func (s coinbaseService) exchangeName() string {
    return coinbaseName
}

func (s coinbaseService) getPairs() []string {
    return coinbasePairs
}

func (s coinbaseService) getTicker(pair string) (error, Ticker) {
    var response coinbaseTicker
    urlP := strings.Replace(pair, "/", "-", -1)
    err := GetJson("https://api.coinbase.com/v2/prices/" + urlP + "/spot", &response)
    return err, response
}

func (t coinbaseTicker) LastPrice() float64 {
    return t.Data.Last
}
