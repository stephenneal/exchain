package exchange

import (
    "fmt"
    "strings"
    "time"
)

type independentReserveService struct{}

type indepResTicker struct {
    DayHighestPrice                  float64   `json:"DayHighestPrice"`
    DayLowestPrice                   float64   `json:"DayLowestPrice"`
    DayAvgPrice                      float64   `json:"DayAvgPrice"`
    DayVolumeXbt                     float64   `json:"DayVolumeXbt"`
    DayVolumeXbtInSecondaryCurrrency float64   `json:"DayVolumeXbtInSecondaryCurrrency"`
    CurrentLowestOfferPrice          float64   `json:"CurrentLowestOfferPrice"`
    CurrentHighestBidPrice           float64   `json:"CurrentHighestBidPrice"`
    Last                             float64   `json:"LastPrice"`
    PrimaryCurrencyCode              string    `json:"PrimaryCurrencyCode"`
    SecondaryCurrencyCode            string    `json:"SecondaryCurrencyCode"`
    CreatedTimestampUtc              time.Time `json:"CreatedTimestampUtc"`
}

var indepResCurr = map[string][]string {
        AUD: {BCH, BTC, ETH},
        USD: {BCH, BTC, ETH},
    }

func (s independentReserveService) getCurrencies() (error, map[string][]string) {
    return nil, indepResCurr
}

func (s independentReserveService) getLastPrice(base string, quot string) (error, float64) {
    var custom indepResTicker
    // For some reason Independent Reserve uses XBT for Bitcoin (BTC)
    p1 := base
    p2 := quot
    if (base == BTC) {
        p1 = "XBT"
    } else if (quot == BTC) {
        p2 = "XBT"
    }
    urlP := fmt.Sprintf("primaryCurrencyCode=%s&secondaryCurrencyCode=%s", strings.ToLower(p1), strings.ToLower(p2))
    err := GetJson("https://api.independentreserve.com/Public/GetMarketSummary?" + urlP, &custom)

    if (err != nil) {
        return err, -1
    }
    return err, custom.Last
}
