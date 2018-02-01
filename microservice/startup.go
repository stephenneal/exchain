package microservice

import (
	"net/http"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	mux := http.NewServeMux()

	mux.Handle("/pub/v1/getTicker", getTickerHandler)
	mux.Handle("/pub/v1/getTickers", getTickersHandler)

	http.Handle("/", accessControl(mux))
	http.Handle("/metrics", promhttp.Handler())

	if !strings.HasPrefix(":", port) {
		port = ":" + port
	}
	errs := make(chan error, 2)
	go func() {
		logger.Log("transport", "http", "address", port, "msg", "listening")
		errs <- http.ListenAndServe(port, nil)
	}()
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Log("terminated", <-errs)
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}