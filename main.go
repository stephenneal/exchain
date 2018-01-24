package main

import (
    //"time"

    "github.com/stephenneal/exchain/api"
)

func main() {
	//ticker := time.NewTicker(time.Millisecond * 1000)
    //go func() {
    //    for range ticker.C {
        for _, p := range api.GetAllPairs() {
            api.RefreshTicker(p)
        }
        api.Derive(api.FIAT_USD, api.FIAT_AUD)
        api.PrintTickers()
    //    }
    //}()
    //time.Sleep(time.Millisecond * 3000)
    //ticker.Stop()
}
