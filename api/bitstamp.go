package api

import (
    "errors"
    "fmt"
    "strconv"
    "strings"

    "github.com/romana/rlog"
)

type bitstampService struct{}

type bitstampTicker struct {
    High      string `json:"high"`
    Last      string `json:"last"`
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

func (s bitstampService) name() string {
    return "Bitstamp"
}

func (s bitstampService) getTicker(pair string) (error, Ticker) {
    var response bitstampTicker
    urlP := strings.ToLower(strings.Replace(pair, "/", "", -1))
    err := GetJson("https://www.bitstamp.net/api/v2/ticker/" + urlP, &response)
    if err != nil {
        return err, nil
    } else if (response.Last == "") {
        return errors.New(fmt.Sprintf("%s (%s); not found", s.name(), pair)), nil
    }
    return nil, response
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

func (t bitstampTicker) LastPrice() float64 {
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
