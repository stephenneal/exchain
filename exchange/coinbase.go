package exchange

import (
    "strings"
)

type coinbaseService struct{}

type coinbaseTicker struct {
    Data struct {
        Base       string `json:"base"`
        Currency   string `json:"currency"`
        Last       float64 `json:"amount,string"`
    } `json:"data"`
}

func (s coinbaseService) getTicker(pair string) (error, SimpleTicker) {
    var custom coinbaseTicker
    urlP := strings.Replace(pair, "/", "-", -1)
    err := GetJson("https://api.coinbase.com/v2/prices/" + urlP + "/spot", &custom)

    var r SimpleTicker
    if (err == nil) {
        r = SimpleTicker { custom.Data.Last }
    }
    return err, r
}
