package microservice

import (
    "github.com/stephenneal/exchain/data"
    "github.com/stephenneal/exchain/exchange"
)

// TickerService provides operations to get ticker info.
type TickerService interface {
    GetTicker(pair string) (error, []data.Ticker)
    GetTickers() (error, []data.Ticker) 
}

type tickerService struct{}

func (tickerService) GetTicker(pair string) (error, []data.Ticker) {
    return exchange.GetTicker(pair)
}

func (tickerService) GetTickers() (error, []data.Ticker) {
    return exchange.GetTickers()
}

// ServiceMiddleware is a chainable behavior modifier for this service.
type ServiceMiddleware func(TickerService) TickerService
