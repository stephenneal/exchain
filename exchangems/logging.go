package exchangems

import (
	"time"

    "github.com/stephenneal/exchain/exchange"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (mw loggingService) GetTicker(s string) (err error, ticker []exchange.Ticker) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "getTicker",
			"input", s,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.Service.GetTicker(s)
}

func (mw loggingService) GetTickers() (error, []exchange.Ticker) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "getTickers",
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.Service.GetTickers()
}
