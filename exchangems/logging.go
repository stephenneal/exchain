package exchangems

import (
	"time"

    "github.com/stephenneal/exchain/exchange"

	kitlog "github.com/go-kit/kit/log"
)

type loggingService struct {
	logger kitlog.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger kitlog.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (mw loggingService) GetTicker(pair string) (err error, ticker []exchange.Ticker) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "getTicker",
			"input", pair,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.Service.GetTicker(pair)
}

func (mw loggingService) GetTickers() (error, []exchange.TickerSummary) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "getTickers",
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.Service.GetTickers()
}
