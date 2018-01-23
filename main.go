package main

import (
    "time"

    "github.com/stephenneal/exchain/api"
)

func main() {
	api.GetRates()

	ticker := time.NewTicker(time.Millisecond * 1000)
    go func() {
        for range ticker.C {
            api.GetTicker(api.ETH_USD)
            api.GetTicker(api.ETH_AUD)
		    //bitstamp.GetTicker("ETH/USD")
		    // api.BitstampTradingPairs()
		    //btcmarkets.GetTicker("ETH/AUD")
        }
    }()
    time.Sleep(time.Millisecond * 3000)
    ticker.Stop()
}
