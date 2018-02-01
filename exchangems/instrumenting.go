package exchangems

import (
	"time"

    "github.com/stephenneal/exchain/exchange"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
	}
}

func (mw instrumentingService) GetTicker(s string) (err error, ticker []exchange.Ticker) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getTicker", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Service.GetTicker(s)
}

func (mw instrumentingService) GetTickers() (error, []exchange.Ticker) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getTickers", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Service.GetTickers()
}
