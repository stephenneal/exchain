package microservice

import (
	"time"

    "github.com/stephenneal/exchain/data"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   TickerService
}

func (mw loggingMiddleware) GetTicker(s string) (err error, ticker []data.Ticker) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "getTicker",
			"input", s,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.GetTicker(s)
}

func (mw loggingMiddleware) GetTickers() (error, []data.Ticker) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "getTickers",
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.GetTickers()
}
