package exchange

import (
    "fmt"
    "strings"
    "time"
)

type indepReserveService struct{}

type indepReserveTicker struct {
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

func (s indepReserveService) getLastPrice(pair TradingPair) (error, float64) {
    var custom indepReserveTicker
    urlP := fmt.Sprintf("primaryCurrencyCode=%s&secondaryCurrencyCode=%s", strings.ToLower(pair.One), strings.ToLower(pair.Two))
    err := GetJson("https://api.independentreserve.com/Public/GetMarketSummary?" + urlP, &custom)

    if (err != nil) {
        return err, -1
    }
    return err, custom.Last
}
