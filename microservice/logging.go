package microservice

import (
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   TickerService
}

func (mw loggingMiddleware) RefreshTickers() {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "refreshTickers",
			"took", time.Since(begin),
		)
	}(time.Now())

	mw.next.RefreshTickers()
	return
}

func (mw loggingMiddleware) RefreshTicker(s string) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "refreshTicker",
			"input", s,
			"took", time.Since(begin),
		)
	}(time.Now())

	mw.next.RefreshTicker(s)
	return
}

func (mw loggingMiddleware) PrintTickers() {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "printTickers",
			"took", time.Since(begin),
		)
	}(time.Now())

	mw.next.PrintTickers()
	return
}
