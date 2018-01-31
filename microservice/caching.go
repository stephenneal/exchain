package microservice

import (
    "github.com/stephenneal/exchain/data"

    "github.com/patrickmn/go-cache"
)

const (
	allKey = "ALL_PAIRS"
)

type cachingMiddleware struct {
	caching *cache.Cache
	next  TickerService
}

func (mw cachingMiddleware) GetTicker(pair string) (error, []data.Ticker) {
	cached, found := mw.caching.Get(pair)
	if (found) {
		return nil, cached.([]data.Ticker)
	}

	err, resp := mw.next.GetTicker(pair)
	if (err != nil) {
		return err, resp
	}
	mw.caching.Set(pair, resp, cache.DefaultExpiration)
	return err, resp
}

func (mw cachingMiddleware) GetTickers() (error, []data.Ticker) {
	cached, found := mw.caching.Get(allKey)
	if (found) {
		return nil, cached.([]data.Ticker)
	}

	err, resp := mw.next.GetTickers()
	if (err != nil) {
		return err, resp
	}
	mw.caching.Set("pair", resp, cache.DefaultExpiration)
	return err, resp
}
