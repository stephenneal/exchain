package api

import (
    "errors"
    "fmt"
    "strconv"
    "strings"

    "github.com/romana/rlog"
)

type binanceService struct{}

type binanceTicker struct {
    Symbol    string `json:"symbol"`
    Last      string `json:"price"`
}

func (s binanceService) name() string {
    return "Binance"
}

func (s binanceService) getTicker(pair string) (error, Ticker) {
    var response binanceTicker
    urlP := strings.Replace(pair, "/", "", -1)
    err := GetJson("https://api.binance.com/api/v3/ticker/price?symbol=" + urlP, &response)
    if err != nil {
        return err, nil
    } else if (response.Last == "") {
        return errors.New(fmt.Sprintf("%s (%s); not found", s.name(), pair)), nil
    }
    return nil, response
}

func (t binanceTicker) LastPrice() float64 {
    if (t.Last == "") {
        return -1
    }
    fPrice, err := strconv.ParseFloat(t.Last, 64)
    if err != nil {
        rlog.Error(err)
        return -1
    }
    return fPrice
}
