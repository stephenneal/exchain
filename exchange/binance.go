package exchange

import (
    "strings"
)

const (
	binanceName = "Binance"
)

var (
    binancePairs = []string{
        BTC_USDT,
        ETH_BTC,
        ETH_USDT,
    }
)

type binanceService struct{}

type binanceTicker struct {
    Symbol    string `json:"symbol"`
    Last      float64 `json:"price,string"`
}

func (s binanceService) getPairs() []string {
    return binancePairs
}

func (s binanceService) exchangeName() string {
    return binanceName
}

func (s binanceService) getTicker(pair string) (error, Ticker) {
    var response binanceTicker
    urlP := strings.Replace(pair, "/", "", -1)
    err := GetJson("https://api.binance.com/api/v3/ticker/price?symbol=" + urlP, &response)
    return err, response
}

func (t binanceTicker) LastPrice() float64 {
    return t.Last
}
