package api

import (
    "fmt"
    "sort"
    "strings"
    "time"

    "github.com/patrickmn/go-cache"
    "github.com/romana/rlog"
)

type Ticker interface {
    LastPrice()    float64
}

type Exchange interface {
    name()            string
    getTicker(string) (error, Ticker)
}

type CacheableTicker struct {
    exchange   string
    pair       string
    lastPrice  float64
    ticker     Ticker
    exchRate   float64
    errorCount int
    lastMod    time.Time
}

const (
    FIAT_AUD = "AUD"
    FIAT_USD = "USD"
    TOK_BTC = "BTC"
    TOK_ETH = "ETH"

    ETH_AUD = TOK_ETH + "/" + FIAT_AUD
    ETH_USD = TOK_ETH + "/" + FIAT_USD
    ETH_BTC = TOK_ETH + "/" + TOK_BTC
)

var (
    bitstamp = bitstampService{}
    btcm = btcmService{}
    binance = binanceService{}
    exByPairs = map[string][]Exchange {
        ETH_AUD : { btcm },
        ETH_USD : { bitstamp },
        ETH_BTC : { binance },
    }

    fiatRates   = cache.New(1*time.Minute, 2*time.Minute)
    tickerCache = cache.New(10*time.Second, 1*time.Minute)
)

// Get ticker for pair from each exchange
func RefreshTicker(pair string) {
    rlog.Debugf("\nGetTicker (%s)", pair)
    var exArr []Exchange
    var ok bool
    if exArr, ok = exByPairs[pair]; ok {
        for _, ex := range exArr {
            // Check the cache
            cacheKey := pair + "-" + ex.name();
            var cached CacheableTicker
            e, found := tickerCache.Get(cacheKey)
            if found {
                cached = e.(CacheableTicker)
                elapsed := int64(time.Now().Sub(cached.lastMod) / time.Millisecond)
                if (elapsed < 10000) {
                    rlog.Debugf("%s (%s): use cached instance; elapsed = %d", ex.name(), pair, elapsed)
                    return
                }
            }

            // Get / update the ticker
            rlog.Debugf("%s (%s): get/update ticker", ex.name(), pair)
            err, ticker := ex.getTicker(pair)
            if err != nil {
                rlog.Error(err)
                // If there is an expired cache entry, just use that
                //if (cached != nil) {
                //    rlog.Errorf("%s (%s): use cached instance", ex.name(), pair)
                //}
            } else {
                rlog.Debugf("%s (%s): %f", ex.name(), pair, ticker.LastPrice())
                tickerCache.Set(cacheKey, CacheableTicker{
                    exchange: ex.name(),
                    pair: pair,
                    exchRate: 0,
                    lastPrice: ticker.LastPrice(),
                    ticker: ticker,
                    errorCount: 0,
                    lastMod: time.Now(),
                }, cache.DefaultExpiration)
            }
        }
    }
}

func GetTickers() {
    // To store the keys in slice in sorted order
    var keys []string
    items := tickerCache.Items()
    for k := range items {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    rlog.Infof("\n")
    for _, k := range keys {
        cached := items[k].Object.(CacheableTicker)
        rlog.Infof(cached.String())
        //rlog.Infof("%s (%s): %f", cached.exchange, cached.pair, cached.ticker.LastPrice())
    }
}

// Derive tickers for currencies not supported by the exchange
func Derive(base string, alt string) {
    err, baseToAltRate := getFiatRate(base, alt)
    if err != nil {
        rlog.Error(err)
        return
    }
    err, altToBaseRate := getFiatRate(alt, base)
    if err != nil {
        rlog.Error(err)
        return
    }

    rlog.Debugf("%s/%s: %f", base, alt, baseToAltRate)
    rlog.Debugf("%s/%s: %f", alt, base, altToBaseRate)
    if (baseToAltRate < 0) {
        return
    }
    if (altToBaseRate < 0) {
        return
    }

    for _, v := range tickerCache.Items() {
        cached := v.Object.(CacheableTicker)
        rlog.Debugf(cached.String())
        splitPair := strings.Split(cached.pair, "/")
        rlog.Debugf("1: %s 2:%s", splitPair[0], splitPair[1])
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
        newCacheKey := newPair + "-" + cached.exchange
        _, found := tickerCache.Get(newCacheKey)
        if found {
            // Already cached, our work here is done
            continue
        }

        price := cached.ticker.LastPrice() * rate
        rlog.Debugf("%s (%s->derived): %f", cached.exchange, newPair, price)
        tickerCache.Set(newCacheKey, CacheableTicker{
            exchange: cached.exchange,
            pair: newPair,
            exchRate: rate,
            lastPrice: price,
            ticker: nil,
            errorCount: 0,
            lastMod: time.Now(),
        }, cache.DefaultExpiration)
    }
}

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

func (t CacheableTicker) String() string {
    var rateStr string
    if (t.exchRate > 0) {
        rateStr = fmt.Sprintf("; exch. rate = %f", t.exchRate)
    }
    return fmt.Sprintf("%s (%s); %f%s", t.pair, t.exchange, t.lastPrice, rateStr)
 }