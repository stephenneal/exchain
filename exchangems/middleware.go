package exchangems

import (
	"time"

	kitlog "github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type ServiceMiddleware func(Service) Service

func InstrumentingMiddleware(counter metrics.Counter, latency metrics.Histogram) ServiceMiddleware {
	return func(next Service) Service {
		return &instrumentingMiddleware{
			requestCount:   counter,
			requestLatency: latency,
			next:           next,
		}
	}
}

func LoggingMiddleware(logger kitlog.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type instrumentingMiddleware struct {
	next	       Service
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
}

type loggingMiddleware struct {
	next   Service
	logger kitlog.Logger
}

// ------------ Instrumenting Middleware ----------------

func (mw instrumentingMiddleware) GetTickers(base string, quot string) (error, []Ticker) {
	defer func(begin time.Time) {
		lvs := []string{"method", "tickers", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.next.GetTickers(base, quot)
}

/*
func (mw instrumentingMiddleware) GetTickers() (error, []TickerSummary) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getTickers", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.next.GetTickers()
}
*/

// ------------ Logging Middleware ----------------

func (mw loggingMiddleware) GetTickers(base string, quot string) (error, []Ticker) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "getTickers",
			"base", base,
			"quot", quot,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.GetTickers(base, quot)
}

/*
func (mw loggingMiddleware) GetTickers() (error, []TickerSummary) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "getTickers",
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.next.GetTickers()
}
*/