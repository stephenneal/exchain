package exchange

import (
    "strings"

    "github.com/romana/rlog"
)

type bitstampService struct{}

type bitstampTicker struct {
    High      string `json:"high"`
    Last      float64 `json:"last,string"`
    Timestamp string `json:"timestamp"`
    Bid       string `json:"bid"`
    Vwap      string `json:"vwap"`
    Volume    string `json:"volume"`
    Low       string `json:"low"`
    Ask       string `json:"ask"`
    Open      string `json:"open"`
}

type tradingPair []struct {
    BaseDecimals    int    `json:"base_decimals"`
    MinimumOrder    string `json:"minimum_order"`
    Name            string `json:"name"`
    CounterDecimals int    `json:"counter_decimals"`
    Trading         string `json:"trading"`
    URLSymbol       string `json:"url_symbol"`
    Description     string `json:"description"`
}

func (s bitstampService) getLastPrice(pair TradingPair) (error, float64) {
    var custom bitstampTicker
    urlP := strings.ToLower(pair.Pair(""))
    err := GetJson("https://www.bitstamp.net/api/v2/ticker/" + urlP, &custom)

    if (err != nil) {
        return err, -1
    }
    return err, custom.Last
}

// FIXME add this to the Exchange API
func GetTradingPairs() {
    var response tradingPair

    if err := GetJson("https://www.bitstamp.net/api/v2/trading-pairs-info/", &response); err != nil {
        rlog.Error(err)
    } else {
        for _, elem := range response {
            rlog.Info(elem.Name)
        }
    }
}
