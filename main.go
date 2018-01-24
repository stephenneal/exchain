package main

import (
    //"time"

    "github.com/stephenneal/exchain/api"
)

func main() {
	//ticker := time.NewTicker(time.Millisecond * 1000)
    //go func() {
    //    for range ticker.C {
            api.RefreshTicker(api.ETH_USD)
            api.RefreshTicker(api.ETH_USDT)
            api.RefreshTicker(api.ETH_AUD)
            api.RefreshTicker(api.ETH_BTC)
            api.Derive(api.FIAT_USD, api.FIAT_AUD)
            api.GetTickers()
    //    }
    //}()
    //time.Sleep(time.Millisecond * 3000)
    //ticker.Stop()
}
