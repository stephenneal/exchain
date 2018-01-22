package bitstamp

import (
    "fmt"
    "strings"

    "github.com/stephenneal/exchain/api"

    "github.com/romana/rlog"
)

const (
    TICKER_URL = "https://www.bitstamp.net/api/v2/ticker/"
	TRADING_PAIRS_URL = "https://www.bitstamp.net/api/v2/trading-pairs-info/"
)

type TickerResponse struct {
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

type TradingPair []struct {
	BaseDecimals    int    `json:"base_decimals"`
	MinimumOrder    string `json:"minimum_order"`
	Name            string `json:"name"`
	CounterDecimals int    `json:"counter_decimals"`
	Trading         string `json:"trading"`
	URLSymbol       string `json:"url_symbol"`
	Description     string `json:"description"`
}

func Ticker(pair string) {
    var response TickerResponse

    p := strings.ToLower(strings.Replace(pair, "/", "", -1))
    if err := api.GetJson(TICKER_URL + p, &response); err != nil {
        rlog.Error(err)
    } else {
        rlog.Info(fmt.Sprintf("Bitstamp (%s); Last=%s; High=%s; Low=%s", pair, response.Last, response.High, response.Low))
    }
}

func TradingPairs() {
    var response TradingPair

    if err := api.GetJson(TRADING_PAIRS_URL, &response); err != nil {
        rlog.Error(err)
    } else {
		for _, elem := range response {
			rlog.Info(elem.Name)
		}
    }
}
