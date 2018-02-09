package main

import (
    "net/http"
    "fmt"
    "os"
    "os/signal"
    "strings"
    "syscall"

    "github.com/stephenneal/exchain/exchangems"

//    "github.com/patrickmn/go-cache"

    stdprometheus "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"

    kitlog "github.com/go-kit/kit/log"
    kitprometheus "github.com/go-kit/kit/metrics/prometheus"
)

func main() {
    var logger kitlog.Logger
    {
        logger = kitlog.NewLogfmtLogger(os.Stderr)
        logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC)
        logger = kitlog.With(logger, "caller", kitlog.DefaultCaller)
    }

    httpAddr := os.Getenv("PORT")
    logger.Log("$PORT", httpAddr)
    if httpAddr == "" {
        logger.Log("message", "$PORT must be set")
        os.Exit(1)
    } else if !strings.HasPrefix(":", httpAddr) {
        httpAddr = ":" + httpAddr
    }


    //caching := cache.New(1*time.Minute, 10*time.Minute)

    fieldKeys := []string{"method", "error"}

    var s exchangems.Service
    {
        s = exchangems.NewService()
        //s = exchangems.NewCachingService(caching, s)
        s = exchangems.LoggingMiddleware(logger)(s)
        s = exchangems.InstrumentingMiddleware(
            kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
                Namespace: "api",
                Subsystem: "exchange_service",
                Name:      "request_count",
                Help:      "Number of requests received.",
            }, fieldKeys),
            kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
                Namespace: "api",
                Subsystem: "exchange_service",
                Name:      "request_latency_microseconds",
                Help:      "Total duration of requests in microseconds.",
            }, fieldKeys))(s)
    }

    var h http.Handler
    {
        h = exchangems.MakeHTTPHandler(s, kitlog.With(logger, "component", "HTTP"))
    }

    mux := http.NewServeMux()
    mux.Handle("/pub/v1/", h)

    http.Handle("/", accessControl(mux))
    http.Handle("/metrics", promhttp.Handler())

    errs := make(chan error, 2)
    go func() {
        c := make(chan os.Signal)
        signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
        errs <- fmt.Errorf("%s", <-c)
    }()
    go func() {
        logger.Log("transport", "http", "address", httpAddr, "msg", "listening")
        errs <- http.ListenAndServe(httpAddr, mux)
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
