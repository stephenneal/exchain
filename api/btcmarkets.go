package api

import (
    "time"

    "github.com/romana/rlog"
)

type btcmService struct{}

type btcmTicker struct {
	BestBid    float64 `json:"bestBid"`
	BestAsk    float64 `json:"bestAsk"`
	Last       float64 `json:"lastPrice"`
	Currency   string  `json:"currency"`
	Instrument string  `json:"instrument"`
	Timestamp  int     `json:"timestamp"`
	Volume24H  float64 `json:"volume24h"`
    errors int
    lastMod   time.Time
}

const (
    BASE_URL = "https://api.btcmarkets.net"
)

func (s btcmService) name() string {
    return "BTC Markets"
}

func (s btcmService) defaultFiat() string {
    return "AUD"
}

func (s btcmService) getTicker(pair string) Ticker {
    var response btcmTicker

    if err := GetJson(BASE_URL + "/market/" + pair + "/tick", &response); err != nil {
        rlog.Error(err)
    } else if (response.Last == 0) {
        rlog.Errorf("%s (%s); not found", s.name(), pair)
    } else {
        rlog.Infof("%s (%s); Last=%f", s.name(), pair, response.LastPrice())
        response.errors = 0
    }
    return response
}

func (t btcmTicker) LastPrice() float64 {
    return t.Last
}

func (t btcmTicker) LastModified() time.Time {
    return t.lastMod
}

func (t btcmTicker) ErrorCount() int {
    return t.errors
}

/*
func tsToTime(sInt int) time.Time {
    msInt := int64(sInt * 1000)
    return time.Unix(0, msInt*int64(time.Millisecond))
}
*/