package exchange

import (
    "fmt"
    "errors"
    "time"

    "github.com/patrickmn/go-cache"
    "github.com/go-kit/kit/log/level"
)

type ExchangeService interface {
    getTicker(string) (error, SimpleTicker)
}

type SimpleTicker struct {
    lastPrice  float64
}

type Ticker2 struct {
    exchange   string
    pair       string
    ticker     *SimpleTicker
    exchRate   float64
    errorCount int
    lastMod    time.Time
}

type TickerProxy interface {
    exchangeName() string
    getPairs() []string
    getTicker(pair string) (err error, ticker Ticker2)
}

// Fetcher implementations
type cachingTicker struct {
    tickerCache *cache.Cache
    next TickerProxy
}

type exchangeTicker struct {
    name string
    tradingPairs []string
    service ExchangeService
}

var (
    allExFetcher = []TickerProxy{
        cachingTicker{ cache.New(5*time.Minute, 10*time.Minute), &exchangeTicker{binanceName, binancePairs, binanceService{}} },
        cachingTicker{ cache.New(5*time.Minute, 10*time.Minute), &exchangeTicker{bitstampName, bitstampPairs, bitstampService{}} },
        cachingTicker{ cache.New(5*time.Minute, 10*time.Minute), &exchangeTicker{btcmName, btcmPairs, btcmService{}} },
        cachingTicker{ cache.New(5*time.Minute, 10*time.Minute), &exchangeTicker{coinbaseName, coinbasePairs, coinbaseService{}} },
    }
)

func GetAllTickers() map[string][]Ticker2 {
    // TODO call exchanges concurrently
    tickers := make(map[string][]Ticker2)
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

func (f cachingTicker) exchangeName() string {
    return f.next.exchangeName()
}

func (f cachingTicker) getPairs() []string {
    return f.next.getPairs()
}

func (f cachingTicker) getTicker(pair string) (error, Ticker2) {
    c, found := tickerCache.Get(pair)
    var t2 Ticker2
    if found {
        t2 := c.(Ticker2)
        return nil, t2
    }
    var err error
    err, t2 = f.next.getTicker(pair)

    if (err == nil) {
        tickerCache.Set(pair, t2, cache.DefaultExpiration)
    }
    return err, t2
}

func (f exchangeTicker) exchangeName() string {
    return f.name
}

func (f exchangeTicker) getPairs() []string {
    return f.tradingPairs
}

func (f exchangeTicker) getTicker(pair string) (error, Ticker2) {
    var t2 Ticker2
    err, t := f.service.getTicker(pair)

    if (err != nil) {
        return err, t2
    }
    if (t.lastPrice <= 0) {
        return errors.New("pair not found"), t2
    }
    t2 = Ticker2{
        exchange: f.name,
        pair: pair,
        exchRate: 0,
        ticker: &t,
        errorCount: 0,
        lastMod: time.Now(),
    }
    return err, t2
}

func (t Ticker2) String() string {
    var rateStr string
    if (t.exchRate > 0) {
        rateStr = fmt.Sprintf("; exch. rate = %f", t.exchRate)
    }
    return fmt.Sprintf("%s (%s); %f%s", t.pair, t.exchange, t.ticker.lastPrice, rateStr)
}
