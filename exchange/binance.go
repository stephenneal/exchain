package exchange

import (
    "strings"
)

type binanceService struct{}

type binanceTicker struct {
    Symbol    string `json:"symbol"`
    Last      float64 `json:"price,string"`
}

func (s binanceService) getTicker(pair string) (error, SimpleTicker) {
    var custom binanceTicker
    urlP := strings.Replace(pair, "/", "", -1)
    err := GetJson("https://api.binance.com/api/v3/ticker/price?symbol=" + urlP, &custom)

    var r SimpleTicker
    if (err == nil) {
    	r = SimpleTicker { custom.Last }
    }
    return err, r
}
