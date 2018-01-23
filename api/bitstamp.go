package api

import (
    "strconv"
    "strings"
    "time"

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
    errors    int
    lastMod   time.Time
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

const (
    TICKER_URL = "https://www.bitstamp.net/api/v2/ticker/"
    TRADING_PAIRS_URL = "https://www.bitstamp.net/api/v2/trading-pairs-info/"
)

var (
    tickerMap = make(map[string]bitstampTicker)
)

func (s bitstampService) name() string {
    return "Bitstamp"
}

func (s bitstampService) defaultFiat() string {
    return "USD"
}

func (s bitstampService) getTicker(pair string) Ticker {
    var response bitstampTicker
    var ok bool
    if response, ok = tickerMap[pair]; ok {
        elapsed := int64(time.Now().Sub(response.LastModified()) / time.Millisecond)
        rlog.Infof("elapsed = %d", elapsed)
        if (elapsed < 2000) {
            rlog.Infof("%s ticker cached (%s); lastMod=%s", s.name(), pair, response.LastModified().String())
            return response
        }
    }

    p := strings.ToLower(strings.Replace(pair, "/", "", -1))
    if err := GetJson(TICKER_URL + p, &response); err != nil {
        rlog.Error(err)
        response.errors = response.errors + 1
        rlog.Errorf("%s (%s); failed. Using cached if it exists", s.name(), pair, response.LastModified().String())
    } else if (response.Last == "") {
        rlog.Errorf("%s (%s); not found", s.name(), pair)
    } else {
        response.errors = 0
        response.lastMod = time.Now()
        tickerMap[pair] = response
        rlog.Infof("%s (%s); Last=%f; High=%s; Low=%s; Time=%s", s.name(), pair, response.LastPrice(), response.High, response.Low, response.LastModified().String())
    }

    return response
}

func GetTradingPairs() {
    var response tradingPair

    if err := GetJson(TRADING_PAIRS_URL, &response); err != nil {
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

func (t bitstampTicker) LastModified() time.Time {
    return t.lastMod
}

func (t bitstampTicker) ErrorCount() int {
    return t.errors
}

/*
func tsToTime(ms string) time.Time {
    sInt, err := strconv.ParseInt(ms, 10, 64)
    if err != nil {
        rlog.Error(err)
        return time.Time{}
    }

    msInt := sInt * 1000
    return time.Unix(0, msInt*int64(time.Millisecond))
}
*/