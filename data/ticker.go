package data

import (
	"fmt"
    "time"
)

type Ticker struct {
    Exchange   string    `json:"exchange"`
    Pair       string    `json:"pair"`
    LastPrice  float64   `json:"last,string"`
    ExchRate   float64   `json:"exchRate,string"`
    ErrorCount int       `json:"errorCount,string"`
    LastMod    time.Time `json:"lastMod,stamp"`
}

func (t Ticker) String() string {
    var rateStr string
    if (t.ExchRate > 0) {
        rateStr = fmt.Sprintf("; exch. rate = %f", t.ExchRate)
    }
    return fmt.Sprintf("%s (%s); %f%s", t.Pair, t.Exchange, t.LastPrice, rateStr)
}
