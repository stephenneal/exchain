package exchange

import (
    "strings"

    //"github.com/romana/rlog"
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

// Pairs supported by the exchange: quote -> base for easier maintainence.
var bitstampCurr = map[string][]string {
    USD: {BCH, BTC, ETH},
}

func (s bitstampService) getPairs() (error, map[string][]string) {
    /* TODO get from Bitstamp and translate
    var response tradingPair
    if err := GetJson("https://www.bitstamp.net/api/v2/trading-pairs-info/", &response); err != nil {
        rlog.Error(err)
    } else {
        for _, elem := range response {
            rlog.Info(elem.Name)
        }
    }
    */
    return nil, bitstampCurr
}

func (s bitstampService) getLastPrice(base string, quot string) (error, float64) {
    var custom bitstampTicker
    urlP := strings.ToLower(base + quot)
    err := GetJson("https://www.bitstamp.net/api/v2/ticker/" + urlP, &custom)

    if (err != nil) {
        return err, -1
    }
    return err, custom.Last
}
