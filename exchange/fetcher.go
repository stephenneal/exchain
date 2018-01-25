package exchange

import (
    "fmt"
    "time"

    "github.com/patrickmn/go-cache"
    "github.com/go-kit/kit/log/level"
)

type Ticker2 struct {
    exchange   string
    pair       string
    lastPrice  float64
    ticker     Ticker
    exchRate   float64
    errorCount int
    lastMod    time.Time
}

type TickerFetcher interface {
    exchangeName() string
    getPairs() []string
    getTicker(pair string) (err error, ticker Ticker2)
}

// Fetcher implementations
type cachingFetcher struct {
    tickerCache *cache.Cache
    next TickerFetcher
}

var (
    allExFetcher = []TickerFetcher{
        cachingFetcher{ cache.New(5*time.Minute, 10*time.Minute), binanceService{} },
        cachingFetcher{ cache.New(5*time.Minute, 10*time.Minute), bitstampService{} },
        cachingFetcher{ cache.New(5*time.Minute, 10*time.Minute), btcmService{} },
        cachingFetcher{ cache.New(5*time.Minute, 10*time.Minute), coinbaseService{} },
    }
)

func GetAllTickers() map[string][]Ticker2 {
    // TODO call exchanges concurrently
    tickers := make(map[string][]Ticker)
    for _, f := range allExFetcher {
        for _, p := range f.getPairs() {
            err, t := f.getTicker(p)
            if (err != nil) {
                level.Error(logger).Log("method", "GetTicker", "pair", p, "exchange", f.exchangeName(), "err", err)
                continue
            }
            tickers[p] = append(tickers[p], t)
        }
    }
    return tickers
}

func GetTicker(pair string) (err error, ticker []Ticker2) {
    // TODO call exchanges concurrently
    var tickers []Ticker2
    for _, f := range allExFetcher {
        var ok bool
        for _, p := range f.getPairs() {
            if p == pair {
                ok = true
                break
            }
        }
        if (!ok) {
            logger.Log("pair", pair, "exchange", f.exchangeName(), "msg", "not supported")
            continue
        }
        err, t := f.getTicker(pair)
        if (err != nil) {
            level.Error(logger).Log("method", "GetTicker", "pair", pair, "exchange", f.exchangeName(), "err", err)
            continue
        }
        tickers = append(tickers, t)
    }
    logger.Log("tickers", fmt.Sprintf("%v", tickers))
    return nil, tickers
}

func (f cachingFetcher) exchangeName() string {
    return f.next.exchangeName()
}

func (f cachingFetcher) getPairs() []string {
    return f.next.getPairs()
}

func (f cachingFetcher) getTicker(pair string) (err error, ticker Ticker2) {
    c, found := tickerCache.Get(pair)
    var t2 Ticker2
    if found {
        t2 := c.(Ticker2)
        return nil, t2
    }
    err, t := f.next.getTicker(pair)

    if (err != nil || t.LastPrice() <= 0) {
        return err, t2
    }
    t2 = Ticker2{
        exchange: f.exchangeName(),
        pair: pair,
        exchRate: 0,
        lastPrice: t.LastPrice(),
        ticker: t,
        errorCount: 0,
        lastMod: time.Now(),
    }, cache.DefaultExpiration)
    return err, t2
}
