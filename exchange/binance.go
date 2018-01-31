package exchange

import (
    "strings"
)

type binanceService struct{}

type binanceTicker struct {
    Symbol    string `json:"symbol"`
    Last      float64 `json:"price,string"`
}

func (s binanceService) getLastPrice(pair string) (error, float64) {
    var custom binanceTicker
    urlP := strings.Replace(pair, "/", "", -1)
    err := GetJson("https://api.binance.com/api/v3/ticker/price?symbol=" + urlP, &custom)

    if (err != nil) {
    	return err, -1
    }
    return err, custom.Last
}
