package microservice

import (
	"net/http"
	"os"
	"time"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
)

func Start(listen *string) {
	logger := log.NewLogfmtLogger(os.Stderr)

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "exchange_group",
		Subsystem: "ticker_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "exchange_group",
		Subsystem: "ticker_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var tickerSvc TickerService
	tickerSvc = tickerService{}
	ticker := time.NewTicker(time.Millisecond * 30000)
    go func() {
       for range ticker.C {
			tickerSvc.RefreshTickers()
       }
    }()
	tickerSvc = loggingMiddleware{logger, tickerSvc}
	tickerSvc = instrumentingMiddleware{requestCount, requestLatency, tickerSvc}

	refreshTickerHandler := httptransport.NewServer(
		makeRefreshTickerEndpoint(tickerSvc),
		decodeRefreshTickerRequest,
		encodeResponse,
	)
	printTickersHandler := httptransport.NewServer(
		makePrintTickersEndpoint(tickerSvc),
		decodeEmptyRequest,
		encodeResponse,
	)

	http.Handle("/refreshTicker", refreshTickerHandler)
	http.Handle("/printTickers", printTickersHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", *listen)
	logger.Log("err", http.ListenAndServe(*listen, nil))
}
