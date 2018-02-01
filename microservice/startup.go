package microservice

import (
	"net/http"
	"os"
	"strings"
	"time"

    "github.com/patrickmn/go-cache"

	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
)

func Start(port string) {
	logger := log.NewLogfmtLogger(os.Stderr)
    caching := cache.New(1*time.Minute, 10*time.Minute)

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
	tickerSvc = cachingMiddleware{caching, tickerSvc}
	tickerSvc = loggingMiddleware{logger, tickerSvc}
	tickerSvc = instrumentingMiddleware{requestCount, requestLatency, tickerSvc}

	getTickerHandler := httptransport.NewServer(
		makeGetTickerEndpoint(tickerSvc),
		decodeTickerRequest,
		httptransport.EncodeJSONResponse,
	)
	getTickersHandler := httptransport.NewServer(
		makeGetTickersEndpoint(tickerSvc),
		decodeEmptyRequest,
		httptransport.EncodeJSONResponse,
	)

	http.Handle("/pub/v1/getTicker", getTickerHandler)
	http.Handle("/pub/v1/getTickers", getTickersHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("msg", "HTTP", "addr", port)
	if !strings.HasPrefix(":", port) {
		port = ":" + port
	}
	logger.Log("err", http.ListenAndServe(port, nil))
}
