package api

import (
    "time"

    //"github.com/romana/rlog"
)

type Ticker interface {
    LastPrice()    float64
    LastModified() time.Time
    ErrorCount()   int
}

type Exchange interface {
    name()            string
    defaultFiat() string
    getTicker(string)       Ticker
}

const (
    EX_BITSTAMP   = "BITSTAMP"
    EX_BTCMARKETS = "BTCMARKETS"

    ETH_AUD = "ETH/AUD"
    ETH_USD = "ETH/USD"
)

var (
    exchangeMap = map[string]Exchange {
        EX_BITSTAMP : bitstampService{},
        EX_BTCMARKETS : btcmService{},
    }
)

func GetTicker(pair string) {
    for _, ex := range exchangeMap {
        ex.getTicker(pair)
    }
}
