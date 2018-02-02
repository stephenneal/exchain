package exchangems

import (
    "github.com/stephenneal/exchain/exchange"

    "github.com/patrickmn/go-cache"
)

const (
	allKey = "ALL_PAIRS"
)

type cachingService struct {
	caching *cache.Cache
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewCachingService(caching *cache.Cache, s Service) Service {
	return &cachingService{caching, s}
}

func (mw cachingService) GetTicker(pair string) (error, []exchange.Ticker) {
	cached, found := mw.caching.Get(pair)
	if (found) {
		return nil, cached.([]exchange.Ticker)
	}

	err, resp := mw.Service.GetTicker(pair)
	if (err != nil) {
		return err, resp
	}
	mw.caching.Set(pair, resp, cache.DefaultExpiration)
	return err, resp
}

func (mw cachingService) GetTickers() (error, []exchange.TickerSummary) {
	cached, found := mw.caching.Get(allKey)
	if (found) {
		return nil, cached.([]exchange.TickerSummary)
	}

	err, resp := mw.Service.GetTickers()
	if (err != nil) {
		return err, resp
	}
	mw.caching.Set(allKey, resp, cache.DefaultExpiration)
	return err, resp
}
