package api

import (
    "strings"
)

type binanceService struct{}

type binanceTicker struct {
    Symbol    string `json:"symbol"`
    Last      float64 `json:"price,string"`
}

func (s binanceService) name() string {
    return "Binance"
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
