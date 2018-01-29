package exchange

import (
    "fmt"
    "os"
    "sort"
    "strings"
    "time"

    "github.com/patrickmn/go-cache"
    "github.com/go-kit/kit/log"
    "github.com/go-kit/kit/log/level"
)

type Exchange struct {
    name string
    tradingPairs []string
    service ExchangeService
}

type CacheableTicker struct {
    exchange   string
    pair       string
    lastPrice  float64
    ticker     *SimpleTicker
    exchRate   float64
    errorCount int
    lastMod    time.Time
}

const (
    FIAT_AUD = "AUD"
    FIAT_USD = "USD"
    TOK_BTC  = "BTC"
    TOK_ETH  = "ETH"
    TOK_USDT = "USDT"

    BTC_AUD  = TOK_BTC + "/" + FIAT_AUD
    BTC_USD  = TOK_BTC + "/" + FIAT_USD
    BTC_USDT = TOK_BTC + "/" + TOK_USDT

    ETH_AUD  = TOK_ETH + "/" + FIAT_AUD
    ETH_BTC  = TOK_ETH + "/" + TOK_BTC
    ETH_USD  = TOK_ETH + "/" + FIAT_USD
    ETH_USDT = TOK_ETH + "/" + TOK_USDT

    binanceName = "Binance"
    bitstampName = "Bitstamp"
    btcmName = "BTCMarkets"
    coinbaseName = "Coinbase"
)

var (
    btcmPairs = []string{
        BTC_AUD,
        ETH_AUD,
    }
    binancePairs = []string{
        BTC_USDT,
        ETH_BTC,
        ETH_USDT,
    }
    bitstampPairs = []string{
        BTC_USD,
        ETH_USD,
    }
    coinbasePairs = []string{
        BTC_AUD,
        BTC_USD,
        ETH_AUD,
        ETH_USD,
    }

    logger = level.NewFilter(log.NewLogfmtLogger(os.Stderr), level.AllowInfo())

    binance = Exchange {binanceName, binancePairs, binanceService{}}
    bitstamp = Exchange {bitstampName, bitstampPairs, bitstampService{}}
    btcm = Exchange {btcmName, btcmPairs, btcmService{}}
    coinbase = Exchange {coinbaseName, coinbasePairs, coinbaseService{}}

    exByPairs = map[string][]Exchange {
        BTC_AUD : { btcm, coinbase },
        BTC_USD : { bitstamp, coinbase },
        BTC_USDT: { binance },
        ETH_AUD : { btcm, coinbase },
        ETH_BTC : { binance },
        ETH_USD : { bitstamp, coinbase },
        ETH_USDT: { binance },
    }
    allPairs = make([]string, len(exByPairs))

    fiatRates   = cache.New(5*time.Minute, 10*time.Minute)
    tickerCache = cache.New(5*time.Minute, 10*time.Minute)
)

func GetAllPairs() []string {
    i := 0
    if (allPairs[i] == "") {
        for k := range exByPairs {
            allPairs[i] = k
            i++
        }
    }
    return allPairs
}

func RefreshTickers() {
    for _, p := range GetAllPairs() {
        RefreshTicker(p)
    }
    Derive(FIAT_USD, FIAT_AUD)
}

// Get ticker for pair from each exchange
func RefreshTicker(pair string) {
    level.Debug(logger).Log("method", "RefreshTicker", "pair", pair)
    var exArr []Exchange
    var ok bool
    if exArr, ok = exByPairs[pair]; ok {
        for _, ex := range exArr {
            // Check the cache
            cacheKey := pair + "-" + ex.name;
            var cached CacheableTicker
            e, found := tickerCache.Get(cacheKey)
            if found {
                cached = e.(CacheableTicker)
                elapsed := int64(time.Now().Sub(cached.lastMod) / time.Millisecond)
                if (elapsed < 10000) {
                    level.Debug(logger).Log("method", "RefreshTicker", "pair", pair, "exchange", ex.name, "cache", "true")
                    return
                }
            }

            // Get / update the ticker
            level.Debug(logger).Log("method", "RefreshTicker", "pair", pair, "exchange", ex.name, "cache", "false")
            err, ticker := ex.service.getTicker(pair)
            if err != nil {
                level.Error(logger).Log("method", "RefreshTicker", "pair", pair, "exchange", ex.name, "message", err)
                // If there is an expired cache entry, just use that
                //if (cached != nil) {
                //    rlog.Errorf("%s (%s): use cached instance", ex.name, pair)
                //}
            } else if (ticker.lastPrice == 0) {
                level.Error(logger).Log("method", "RefreshTicker", "pair", pair, "exchange", ex.name, "message", "empty")
            } else {
                level.Debug(logger).Log("method", "RefreshTicker", "pair", pair, "exchange", ex.name, "lastPrice", ticker.lastPrice)
                tickerCache.Set(cacheKey, CacheableTicker{
                    exchange: ex.name,
                    pair: pair,
                    exchRate: 0,
                    lastPrice: ticker.lastPrice,
                    ticker: &ticker,
                    errorCount: 0,
                    lastMod: time.Now(),
                }, cache.DefaultExpiration)
            }
        }
    }
}

func GetTickers() map[string][]CacheableTicker {
    tickers := make(map[string][]CacheableTicker)
    for _, v := range tickerCache.Items() {
        cached := v.Object.(CacheableTicker)
        tickers[cached.pair] = append(tickers[cached.pair], cached)
    }
    return tickers
}

func PrintTickers() {
    // Sort the keys
    var keys []string
    items := tickerCache.Items()
    for k := range items {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    var prevPair string
    for _, k := range keys {
        cached := items[k].Object.(CacheableTicker)
        // Print a new line in between pairs
        var customStr string
        if (prevPair != cached.pair) {
            // Put USDT with USD but demarcate
            if (strings.HasSuffix(prevPair, FIAT_USD) && strings.HasSuffix(cached.pair, TOK_USDT)) {
                customStr = " (" + TOK_USDT + ")"
            } else {
                prevPair = cached.pair
                logger.Log("PrintTickers", cached.pair)
            }
        }
        if (cached.exchRate > 0) {
            customStr = fmt.Sprintf(" (exch=%f)", cached.exchRate)
        }
        logger.Log("exchange", cached.exchange, "lastPrice", cached.lastPrice, "extra", customStr)
    }
}

// Derive tickers for currencies not supported by the exchange
func Derive(base string, alt string) {
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
        cached := v.Object.(CacheableTicker)
        splitPair := strings.Split(cached.pair, "/")
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

        price := cached.ticker.lastPrice * rate
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
