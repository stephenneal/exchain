package exchange

import (
	"fmt"
    "time"
)

type TradingPair struct {
    One string    `json:"one"`
    Two string    `json:"two"`
}

type Ticker struct {
    Exchange    string      `json:"exchange"`
    Pair        TradingPair `json:"pair, string"`
    LastPrice   float64     `json:"last,string"`
    ExchRate    float64     `json:"exchRate,string"`
    Err         string      `json:"error"`
    LastMod     time.Time   `json:"lastMod,stamp"`
}

func (t Ticker) String() string {
    if (len(t.Err) > 0) {
        return fmt.Sprintf("%s (%s); %s", t.Err)
    }
    var rateStr string
    if (t.ExchRate > 0) {
        rateStr = fmt.Sprintf("; exch. rate = %f", t.ExchRate)
    }
    return fmt.Sprintf("%s (%s); %f%s", t.Pair, t.Exchange, t.LastPrice, rateStr)
}

func (t TradingPair) Pair(separator string) string {
    return fmt.Sprintf("%s%s%s", t.One, separator, t.Two)
}

func (t TradingPair) String() string {
    return fmt.Sprintf("%s/%s", t.One, t.Two)
}
