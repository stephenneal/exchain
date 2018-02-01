package microservice

import (
    "github.com/stephenneal/exchain/exchange"
)

// TickerService provides operations to get ticker info.
type TickerService interface {
    GetTicker(pair string) (error, []exchange.Ticker)
    GetTickers() (error, []exchange.Ticker) 
}

type tickerService struct{}

var ex = exchange.Exchange{}

func (tickerService) GetTicker(pair string) (error, []exchange.Ticker) {
    return ex.GetTicker(pair)
}

func (tickerService) GetTickers() (error, []exchange.Ticker) {
    return ex.GetTickers()
}

// ServiceMiddleware is a chainable behavior modifier for this service.
type ServiceMiddleware func(TickerService) TickerService
