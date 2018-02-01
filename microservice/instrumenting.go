package microservice

import (
	"time"

    "github.com/stephenneal/exchain/exchange"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           TickerService
}

func (mw instrumentingMiddleware) GetTicker(s string) (err error, ticker []exchange.Ticker) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getTicker", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.next.GetTicker(s)
}

func (mw instrumentingMiddleware) GetTickers() (error, []exchange.Ticker) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getTickers", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.next.GetTickers()
}
