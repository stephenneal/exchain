package exchangems

import (
    "github.com/stephenneal/exchain/exchange"
)

// Service is an interface that provides exchange operations.
type Service interface {
    GetTicker(pair string) (error, []exchange.Ticker)
    GetTickers() (error, []exchange.Ticker) 
}

type service struct{}

var ex = exchange.Exchange{}

func (service) GetTicker(pair string) (error, []exchange.Ticker) {
    return ex.GetTicker(pair)
}

func (service) GetTickers() (error, []exchange.Ticker) {
    return ex.GetTickers()
}

// ServiceMiddleware is a chainable behavior modifier for this service.
type ServiceMiddleware func(Service) Service

// NewService creates an exchange service with necessary dependencies.
func NewService() Service {
	return &service{}
}