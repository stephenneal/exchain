package exchangems

import (
    //"context"
    //"errors"
	"fmt"
    "time"
)

// Service is an interface that provides exchange operations.
type Service interface {
    GetTickers(base string, quot string) (error, []Ticker)
    //GetTickers() (error, []TickerSummary) 
}

type Ticker struct {
    Exchange    string      `json:"exchange"`
    Base string             `json:"base"`
    Quot string             `json:"quot"`
    LastPrice   float64     `json:"last,string"`
    ExchRate    float64     `json:"exchRate,string"`
    Err         string      `json:"error"`
    LastMod     time.Time   `json:"lastMod,stamp"`
}

/*
type TickerSummary struct {
    Pair         string    `json:"pair, string"`
    HighestPrice float64   `json:"highPrice,string"`
    LowestPrice  float64   `json:"lowPrice,string"`
    Err          string    `json:"error"`
    LastMod      time.Time `json:"lastMod,stamp"`
    Tickers      []Ticker  `json:"tickers"`
}
*/

//var (
//	ErrNotFound        = errors.New("not found")
//)

type service struct{}

func NewService() Service {
	return &service{}
}

func (s service) GetTickers(base string, quot string) (error, []Ticker) {
    var t []Ticker
    return nil, t
}

/*
func (s service) GetTickerSummary() (error, []TickerSummary) {
    var ts []TickerSummary
    return nil, ts
}
*/

func (t Ticker) String() string {
    if (len(t.Err) > 0) {
        return fmt.Sprintf("%s/%s (%s); %s", t.Base, t.Quot, t.Exchange, t.Err)
    }
    var rateStr string
    if (t.ExchRate > 0) {
        rateStr = fmt.Sprintf("; exch. rate = %f", t.ExchRate)
    }
    return fmt.Sprintf("%s/%s (%s); %f%s", t.Base, t.Quot, t.Exchange, t.LastPrice, rateStr)
}

/*
func (t TradingPair) Pair(separator string) string {
    return fmt.Sprintf("%s%s%s", t.Base, separator, t.Quot)
}

func (t TradingPair) String() string {
    return t.Pair("/")
}

func (t TickerSummary) String() string {
    if (len(t.Err) > 0) {
        return fmt.Sprintf("%s; %s", t.Pair, t.Err)
    }
    return fmt.Sprintf("%s; lowest=%f, highest=%f, %s", t.Pair, t.LowestPrice, t.HighestPrice)
}
*/