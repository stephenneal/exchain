package exchange

import (
    "fmt"
    "os"
    "sort"
    "strings"
    "time"

    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/log/level"
    "github.com/patrickmn/go-cache"
)

type exchangeService interface {
    getLastPrice(TradingPair) (error, float64)
}

type Exchange struct {
    name string
    service exchangeService
}

const (
    FIAT_AUD = "AUD"
    FIAT_USD = "USD"

    TOK_BCH  = "BCH"
    TOK_BTC  = "BTC"
    TOK_ETH  = "ETH"
    TOK_USDT = "USDT"

    sep = "-"
    BCH_AUD  = TOK_BCH + sep + FIAT_AUD
    BCH_USD  = TOK_BCH + sep + FIAT_USD

    BTC_AUD  = TOK_BTC + sep + FIAT_AUD
    BTC_USD  = TOK_BTC + sep + FIAT_USD
    BTC_USDT = TOK_BTC + sep + TOK_USDT

    ETH_AUD  = TOK_ETH + sep + FIAT_AUD
    ETH_BTC  = TOK_ETH + sep + TOK_BTC
    ETH_USD  = TOK_ETH + sep + FIAT_USD
    ETH_USDT = TOK_ETH + sep + TOK_USDT
)

var (
    logger = level.NewFilter(log.NewLogfmtLogger(os.Stderr), level.AllowInfo())

    binance  = Exchange {"Binance", binanceService{}}
    bitstamp = Exchange {"Bitstamp", bitstampService{}}
    btcm     = Exchange {"BTCMarkets", btcmService{}}
    coinbase = Exchange {"Coinbase", coinbaseService{}}
    indepReserve = Exchange {"Independent Reserve", indepReserveService{}}

    exByPairs = map[string][]Exchange {
        BCH_AUD : { btcm, indepReserve },
        BCH_USD : { bitstamp, coinbase, indepReserve },
        BTC_AUD : { btcm, coinbase },
        BTC_USD : { bitstamp, coinbase },
        BTC_USDT: { binance },
        ETH_AUD : { btcm, coinbase, indepReserve },
        ETH_BTC : { binance },
        ETH_USD : { bitstamp, coinbase, indepReserve },
        ETH_USDT: { binance },
    }
    allPairs = make([]string, 0, len(exByPairs))

    fiatRates = cache.New(5*time.Minute, 10*time.Minute)
)

func (ex Exchange) AllPairs() []string {
    i := 0
    if (len(allPairs) == 0) {
        for k := range exByPairs {
            allPairs = append(allPairs, k)
            i++
        }
        sort.Strings(allPairs)
    }
    return allPairs
}

func (ex Exchange) GetTickers() (error, []Ticker) {
    // TODO call exchanges concurrently
    var tickers []Ticker
    for _, pair := range ex.AllPairs() {
        err, t := ex.GetTicker(pair)
        if (err != nil) {
            level.Error(logger).Log("method", "GetTickers", "pair", pair, "err", err)
            continue
        }
        tickers = append(tickers, t...)
    }
    level.Debug(logger).Log("method", "GetTickers", "tickers", fmt.Sprintf("%v", tickers))
    return nil, tickers
}

func (ex Exchange) GetTicker(pair string) (error, []Ticker) {
    // TODO call exchanges concurrently
    // TODO deal with timeouts (exchange unavailable)...
    var tickers []Ticker
    if exArr, ok := exByPairs[pair]; ok {
        for _, ex := range exArr {
            splitPair := strings.Split(pair, sep)
            tp := TradingPair { splitPair[0], splitPair[1] }
            err, lastPrice := ex.service.getLastPrice(tp)

            t := Ticker{
                Exchange: ex.name,
                Pair: tp,
                LastMod: time.Now(),
            }

            if (err != nil) {
                t.Err = err.Error()
            } else if (lastPrice <= 0) {
                t.Err = "not found on this exchange"
            } else {
                t.ExchRate = 0
                t.LastPrice = lastPrice
            }
            level.Debug(logger).Log("method", "GetTicker", "ticker", t)
            tickers = append(tickers, t)
        }
    }
    level.Debug(logger).Log("method", "GetTicker", "pair", pair, "tickers", fmt.Sprintf("%v", tickers))
    return nil, tickers
}

/*
func refreshTickers() {
    for _, p := range GetAllPairs() {
        RefreshTicker(p)
    }
    derive(FIAT_USD, FIAT_AUD)
}

func printTickers() {
    // Sort the keys
    var keys []string
    items := tickerCache.Items()
    for k := range items {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    var prevPair string
    for _, k := range keys {
        cached := items[k].Object.(Ticker)
        // Print a new line in between pairs
        var customStr string
        if (prevPair != cached.Pair) {
            // Put USDT with USD but demarcate
            if (strings.HasSuffix(prevPair, FIAT_USD) && strings.HasSuffix(cached.Pair, TOK_USDT)) {
                customStr = " (" + TOK_USDT + ")"
            } else {
                prevPair = cached.Pair
                logger.Log("PrintTickers", cached.Pair)
            }
        }
        if (cached.ExchRate > 0) {
            customStr = fmt.Sprintf(" (exch=%f)", cached.ExchRate)
        }
        logger.Log("exchange", cached.Exchange, "lastPrice", cached.LastPrice, "extra", customStr)
    }
}

// Derive tickers for currencies not supported by the exchange
func derive(base string, alt string) {
    err, baseToAltRate := getFiatRate(base, alt)
    if err != nil {
        level.Error(logger).Log("method", "Derive", "base", base, "alt", alt, "error", err)
        return
    }
    err, altToBaseRate := getFiatRate(alt, base)
    if err != nil {
        level.Error(logger).Log("method", "Derive", "base", base, "alt", alt, "error", err)
        return
    }

    level.Debug(logger).Log("method", "Derive", base + "/" + alt, "rate", baseToAltRate, alt + "/" + base, "rate", altToBaseRate)
    if (baseToAltRate < 0) {
        return
    }
    if (altToBaseRate < 0) {
        return
    }

    for _, v := range tickerCache.Items() {
        cached := v.Object.(Ticker)
        splitPair := strings.Split(cached.Pair, "/")
        fiat := splitPair[1]

        var other string
        var rate float64
        if (fiat == base) {
            other = alt
            rate = baseToAltRate
        } else if (fiat == alt) {
            other = base
            rate = altToBaseRate
        } else {
            // Not deriving this currency
            continue
        }

        // Is the other pair cached already?
        newPair := splitPair[0] + "/" + other
        newCacheKey := newPair + "-" + cached.Exchange
        _, found := tickerCache.Get(newCacheKey)
        if found {
            // Already cached, our work here is done
            continue
        }

        price := cached.LastPrice * rate
        tickerCache.Set(newCacheKey, Ticker{
            Exchange: cached.Exchange,
            Pair: newPair,
            ExchRate: rate,
            LastPrice: price,
            ErrorCount: 0,
            LastMod: time.Now(),
        }, cache.DefaultExpiration)
    }
}
*/

func getFiatRate(base string, alt string) (error, float64) {
    var err error
    var rates Rates
    e, found := fiatRates.Get(base)
    if found {
        rates = e.(Rates)
    } else {
        err, rates = GetFiatRates(base)
        if err != nil {
            return err, -1
        }
        fiatRates.Set(base, rates, cache.DefaultExpiration)
    }
    var rate float64
    switch alt {
    case FIAT_AUD:
        rate = rates.Currencies.AUD
    case FIAT_USD:
        rate = rates.Currencies.USD
    default:
        // FIXME return error
        rate = -1
    }
    return err, rate
}
