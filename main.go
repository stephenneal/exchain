package main

import (
    "time"

    "github.com/stephenneal/exchain/api/bitstamp"
    "github.com/stephenneal/exchain/api/btcmarkets"
    "github.com/stephenneal/exchain/api/fiat"
)

func main() {
	fiat.GetRates()

	ticker := time.NewTicker(time.Millisecond * 1000)
    go func() {
        for range ticker.C {
		    bitstamp.Ticker("ETH/USD")
		    // api.BitstampTradingPairs()
		    btcmarkets.Ticker("ETH/AUD")
        }
    }()
    time.Sleep(time.Millisecond * 3000)
    ticker.Stop()
}
