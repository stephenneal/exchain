package exchangems

import (
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

func (mw cachingService) GetTickers(base string, quot string) (error, []Ticker) {
	cached, found := mw.caching.Get(base)
	if (found) {
		return nil, cached.([]Ticker)
	}

	err, resp := mw.Service.GetTickers(base, quot)
	if (err != nil) {
		return err, resp
	}
	mw.caching.Set(base, resp, cache.DefaultExpiration)
	return err, resp
}

/*
func (mw cachingService) GetTickers() (error, []TickerSummary) {
	cached, found := mw.caching.Get(allKey)
	if (found) {
		return nil, cached.([]TickerSummary)
	}

	err, resp := mw.Service.GetTickers()
	if (err != nil) {
		return err, resp
	}
	mw.caching.Set(allKey, resp, cache.DefaultExpiration)
	return err, resp
}
*/
