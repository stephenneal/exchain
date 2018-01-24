package microservice

import (
    "fmt"
    "os"

    "github.com/stephenneal/exchain/exchange"

    "github.com/go-kit/kit/log"
)

// TickerService provides operations to get ticker info.
type TickerService interface {
    RefreshTickers()
	RefreshTicker(pair string)
	PrintTickers()
}

type tickerService struct{}

func (tickerService) RefreshTickers() {
    exchange.RefreshTickers()
}

func (tickerService) RefreshTicker(pair string) {
	exchange.RefreshTicker(pair)
}

func (tickerService) PrintTickers() {
    logger := log.NewLogfmtLogger(os.Stderr)
    for k, v := range exchange.GetTickers() {
    	logger.Log("pair", k, "tickers", fmt.Sprintf("%v", v))
    }
	//exchange.PrintTickers()
}

// ServiceMiddleware is a chainable behavior modifier for this service.
type ServiceMiddleware func(TickerService) TickerService
