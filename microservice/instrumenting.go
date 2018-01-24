package microservice

import (
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           TickerService
}

func (mw instrumentingMiddleware) RefreshTickers() {
	defer func(begin time.Time) {
		lvs := []string{"method", "refreshTickers", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	mw.next.RefreshTickers()
	return
}

func (mw instrumentingMiddleware) RefreshTicker(s string) {
	defer func(begin time.Time) {
		lvs := []string{"method", "refreshTicker", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	mw.next.RefreshTicker(s)
	return
}

func (mw instrumentingMiddleware) PrintTickers() {
	defer func(begin time.Time) {
		lvs := []string{"method", "printTickers", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	mw.next.PrintTickers()
	return
}