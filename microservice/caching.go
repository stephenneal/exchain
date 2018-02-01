package microservice

import (
    "github.com/stephenneal/exchain/exchange"

    "github.com/patrickmn/go-cache"
)

const (
	allKey = "ALL_PAIRS"
)

type cachingMiddleware struct {
	caching *cache.Cache
	next  TickerService
}

func (mw cachingMiddleware) GetTicker(pair string) (error, []exchange.Ticker) {
	cached, found := mw.caching.Get(pair)
	if (found) {
		return nil, cached.([]exchange.Ticker)
	}

	err, resp := mw.next.GetTicker(pair)
	if (err != nil) {
		return err, resp
	}
	mw.caching.Set(pair, resp, cache.DefaultExpiration)
	return err, resp
}

func (mw cachingMiddleware) GetTickers() (error, []exchange.Ticker) {
	cached, found := mw.caching.Get(allKey)
	if (found) {
		return nil, cached.([]exchange.Ticker)
	}

	err, resp := mw.next.GetTickers()
	if (err != nil) {
		return err, resp
	}
	mw.caching.Set("pair", resp, cache.DefaultExpiration)
	return err, resp
}
